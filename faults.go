package sectorstorage

import (
	"context"
	"os"
	"path/filepath"

	"golang.org/x/xerrors"

	"github.com/filecoin-project/sector-storage/stores"
	"github.com/filecoin-project/specs-actors/actors/abi"
)

// TODO: Track things more actively
type FaultTracker interface {
	CheckProvable(ctx context.Context, spt abi.RegisteredSealProof, sectors []abi.SectorID) ([]abi.SectorID, error)
}

// Returns unprovable sectors
func (m *Manager) CheckProvable(ctx context.Context, spt abi.RegisteredSealProof, sectors []abi.SectorID) ([]abi.SectorID, error) {
	var bad []abi.SectorID

	ssize, err := spt.SectorSize()
	if err != nil {
		return nil, err
	}

	// TODO: More better checks
	for _, sector := range sectors {
		err := func() error {
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			locked, err := m.index.StorageTryLock(ctx, sector, stores.FTSealed|stores.FTCache, stores.FTNone)
			if err != nil {
				return xerrors.Errorf("acquiring sector lock: %w", err)
			}

			if !locked {
				log.Warnw("CheckProvable Sector FAULT: can't acquire read lock", "sector", sector, "sealed")
				bad = append(bad, sector)
				return nil
			}

			lp, _, err := m.localStore.AcquireSector(ctx, sector, spt, stores.FTSealed|stores.FTCache, stores.FTNone, false, stores.AcquireMove)
			if err != nil {
				return xerrors.Errorf("acquire sector in checkProvable: %w", err)
			}

			if lp.Sealed == "" || lp.Cache == "" {
				log.Warnw("CheckProvable Sector FAULT: cache an/or sealed paths not found", "sector", sector, "sealed", lp.Sealed, "cache", lp.Cache)
				bad = append(bad, sector)
				return nil
			}

			toCheck := map[string]int64{
				lp.Sealed:                        1,
				filepath.Join(lp.Cache, "t_aux"): 0,
				filepath.Join(lp.Cache, "p_aux"): 0,
				filepath.Join(lp.Cache, "sc-02-data-tree-r-last.dat"): 0,
			}

			for p, sz := range toCheck {
				st, err := os.Stat(p)
				if err != nil {
					log.Warnw("CheckProvable Sector FAULT: sector file stat error", "sector", sector, "sealed", lp.Sealed, "cache", lp.Cache, "file", p)
					bad = append(bad, sector)
					return nil
				}

				if sz != 0 {
					if st.Size() != int64(ssize)*sz {
						log.Warnw("CheckProvable Sector FAULT: sector file is wrong size", "sector", sector, "sealed", lp.Sealed, "cache", lp.Cache, "file", p, "size", st.Size(), "expectSize", int64(ssize)*sz)
						bad = append(bad, sector)
						return nil
					}
				}
			}

			return nil
		}()
		if err != nil {
			return nil, err
		}
	}

	return bad, nil
}

var _ FaultTracker = &Manager{}
