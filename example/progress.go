package example

import (
	"fmt"
	"net/http"
	"time"

	"sync"

	"github.com/lyismydg/fxgos/service"
)

type ProgressExample struct {
}

func (pe *ProgressExample) ServiceProgress(w http.ResponseWriter, r *http.Request) {
	params := service.GetRequestVars(r, "progressKey")
	progress := service.GetProgress(params["progressKey"], "progressDemo")

	for i := 1; i <= 20; i++ {
		progress.Active(i, fmt.Sprintf("Main Progress %d%%", i))
		time.Sleep(1 * time.Second)
	}

	time.Sleep(5 * time.Second)

	goProgress("SubProgress", 5, 1*time.Second, progress.NewSubProgress(50))

	time.Sleep(5 * time.Second)

	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			goProgress(fmt.Sprintf("Routine00%d", index), 10, 1*time.Second, progress.NewSubProgress(10))
		}(i)
	}

	wg.Wait()

	time.Sleep(5 * time.Second)

	service.ResponseJSON(w, nil, map[string]interface{}{"ProgressDone": true}, http.StatusOK)

}

func goProgress(code string, step int, duration time.Duration, progressSubscribers ...service.ProgressSubscriber) {
	pd := service.NewProgressDispatcher(code, progressSubscribers...)
	defer pd.Success()
	ticker := time.NewTicker(duration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if pd.GetPercent()+step > 100 {
				return
			} else {
				pd.Step(step, fmt.Sprintf("Progress %s = %d", code, pd.GetPercent()+step))
			}
		default:

		}
	}
}

func init() {
	pe := &ProgressExample{}
	myRouter.Router().Path("/progress").HandlerFunc(pe.ServiceProgress)
}
