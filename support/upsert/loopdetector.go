package upsert

import (
	"fmt"
	"sync"

	"github.com/bombsimon/logrusr"
	"github.com/go-logr/logr"
	"github.com/google/go-cmp/cmp"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func newUpdateLoopDetector() *updateLoopDetector {
	return &updateLoopDetector{
		hasNoOpUpdate:    sets.String{},
		updateEventCount: map[string]int{},
		log:              logrusr.NewLogger(func() logrus.FieldLogger { l := logrus.New(); l.SetFormatter(&logrus.JSONFormatter{}); return l }()),
	}
}

// LoopDetectorWarningMessage is logged whenever we detect multiple updates of the same object
// without observering a no-op update.
const LoopDetectorWarningMessage = "WARNING: Object got updated more than one time without a no-op update, this indicates hypershift incorrectly reverting defaulted values"

// If an object got updated more than once a no-op update, we assume it is a bug in our
// code. This is a heuristic that currently happens to work out but might need adjustment
// in the future.
// Once we did a no-op update, we will ignore the object because we assume that if we have
// a bug in the defaulting, we will end up always updating.
const updateLoopThreshold = 2

type updateLoopDetector struct {
	hasNoOpUpdate    sets.String
	lock             sync.RWMutex
	updateEventCount map[string]int
	log              logr.Logger
}

func (*updateLoopDetector) keyFor(obj runtime.Object, key crclient.ObjectKey) string {
	return fmt.Sprintf("%T %s", obj, key.String())
}

func (uld *updateLoopDetector) recordNoOpUpdate(obj crclient.Object, key crclient.ObjectKey) {
	uld.lock.Lock()
	defer uld.lock.Unlock()
	uld.hasNoOpUpdate.Insert(uld.keyFor(obj, key))
}

func (uld *updateLoopDetector) recordActualUpdate(original, modified runtime.Object, key crclient.ObjectKey) {
	cacheKey := uld.keyFor(original, key)
	uld.lock.RLock()
	hasNoOpUpdate := uld.hasNoOpUpdate.Has(cacheKey)
	uld.lock.RUnlock()

	if hasNoOpUpdate {
		return
	}

	uld.lock.Lock()
	uld.updateEventCount[cacheKey]++
	updateEventCount := uld.updateEventCount[cacheKey]
	uld.lock.Unlock()

	if updateEventCount < updateLoopThreshold {
		return
	}

	diff := cmp.Diff(original, modified)
	semanticDeepEqual := equality.Semantic.DeepEqual(original, modified)
	uld.log.Info(LoopDetectorWarningMessage, "type", fmt.Sprintf("%T", modified), "name", key.String(), "diff", diff, "semanticDeepEqual", semanticDeepEqual, "updateCount", updateEventCount)
}
