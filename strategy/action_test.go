// Copyright (c) 2021-2026 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package strategy_test

import (
	"testing"

	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/strategy"
)

func TestAnnotation(t *testing.T) {
	actions := []strategy.Action{strategy.Hold, strategy.Buy, strategy.Sell}
	annotations := []string{"", "B", "S"}

	for i, action := range actions {
		actual := action.Annotation()
		expected := annotations[i]

		if actual != expected {
			t.Fatalf("actual %s expected %s", actual, expected)
		}
	}
}

func TestActionsToAnnotations(t *testing.T) {
	actions := helper.SliceToChan([]strategy.Action{strategy.Hold, strategy.Buy, strategy.Sell})
	expected := helper.SliceToChan([]string{"", "B", "S"})

	actual := strategy.ActionsToAnnotations(actions)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNormalizeActions(t *testing.T) {
	actions := helper.SliceToChan([]strategy.Action{
		strategy.Hold, strategy.Sell, strategy.Sell, strategy.Buy, strategy.Hold,
		strategy.Buy, strategy.Buy, strategy.Sell, strategy.Sell, strategy.Buy,
	})

	expected := helper.SliceToChan([]strategy.Action{
		strategy.Hold, strategy.Hold, strategy.Hold, strategy.Buy, strategy.Hold,
		strategy.Hold, strategy.Hold, strategy.Sell, strategy.Hold, strategy.Buy,
	})

	actual := strategy.NormalizeActions(actions)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDenormalizeActions(t *testing.T) {
	actions := helper.SliceToChan([]strategy.Action{
		strategy.Hold, strategy.Hold, strategy.Hold, strategy.Buy, strategy.Hold,
		strategy.Hold, strategy.Hold, strategy.Sell, strategy.Hold, strategy.Buy,
	})

	expected := helper.SliceToChan([]strategy.Action{
		strategy.Hold, strategy.Hold, strategy.Hold, strategy.Buy, strategy.Buy,
		strategy.Buy, strategy.Buy, strategy.Sell, strategy.Sell, strategy.Buy,
	})

	actual := strategy.DenormalizeActions(actions)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCountActions(t *testing.T) {
	chan1 := helper.SliceToChan[strategy.Action]([]strategy.Action{strategy.Buy})
	chan2 := helper.SliceToChan[strategy.Action]([]strategy.Action{strategy.Hold})
	chan3 := helper.SliceToChan[strategy.Action]([]strategy.Action{strategy.Sell})

	expectedBuy := 1
	expectedHold := 1
	expectedSell := 1

	actualBuy, actualHold, actualSell, ok := strategy.CountActions([]<-chan strategy.Action{
		chan1, chan2, chan3,
	})

	if !ok {
		t.Fatal("not ok")
	}

	if actualBuy != expectedBuy {
		t.Fatalf("actual %v expected %v", actualBuy, expectedBuy)
	}

	if actualHold != expectedHold {
		t.Fatalf("actual %v expected %v", actualHold, expectedHold)
	}

	if actualSell != expectedSell {
		t.Fatalf("actual %v expected %v", actualSell, expectedSell)
	}
}

func TestCountActionsEmpty(t *testing.T) {
	chan1 := helper.SliceToChan[strategy.Action]([]strategy.Action{})
	chan2 := helper.SliceToChan[strategy.Action]([]strategy.Action{})
	chan3 := helper.SliceToChan[strategy.Action]([]strategy.Action{})

	_, _, _, ok := strategy.CountActions([]<-chan strategy.Action{
		chan1, chan2, chan3,
	})

	if ok {
		t.Fatal("is ok")
	}
}

func TestCountTransactions(t *testing.T) {
	actions := helper.SliceToChan([]strategy.Action{strategy.Hold, strategy.Buy, strategy.Hold, strategy.Sell})
	expected := helper.SliceToChan([]int{0, 1, 1, 2})

	actual := strategy.CountTransactions(actions)

	err := helper.CheckEquals(actual, expected)
	if err != nil {
		t.Fatal(err)
	}
}
