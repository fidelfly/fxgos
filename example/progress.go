package example

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/fidelfly/fxgo"
	"github.com/fidelfly/fxgo/httprxr"
)

type ProgressExample struct {
}

func (pe *ProgressExample) ServiceProgress(w http.ResponseWriter, r *http.Request) {
	params := httprxr.GetRequestVars(r, "progressKey")
	progress := fxgo.GetProgress(params["progressKey"], "progressDemo")

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

	time.Sleep(1 * time.Second)

	httprxr.ResponseJSON(w, http.StatusOK, map[string]interface{}{"ProgressDone": true})

}

func goProgress(code string, step int, duration time.Duration, progressSubscribers ...httprxr.ProgressSubscriber) {
	pd := httprxr.NewProgressDispatcher(code, progressSubscribers...)
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
	fxgo.Router().Path("/example/progress").HandlerFunc(pe.ServiceProgress)
}
