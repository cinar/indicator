package momentum

import (
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

// PringsSpecialK implements Martin Pring's Special K momentum indicator.
// It composes multiple Rate-of-Change (ROC) series smoothed by Simple Moving Averages (SMA)
// and outputs a weighted sum aligned to the slowest path so all terms are time-synchronized.
// See Compute for the exact composition and weights.
type PringsSpecialK[T helper.Float] struct {
	Roc10  *trend.Roc[T]
	Roc15  *trend.Roc[T]
	Roc20  *trend.Roc[T]
	Roc30  *trend.Roc[T]
	Roc40  *trend.Roc[T]
	Roc65  *trend.Roc[T]
	Roc75  *trend.Roc[T]
	Roc100 *trend.Roc[T]
	Roc195 *trend.Roc[T]
	Roc265 *trend.Roc[T]
	Roc390 *trend.Roc[T]
	Roc530 *trend.Roc[T]

	Sma10Roc10   *trend.Sma[T]
	Sma10Roc15   *trend.Sma[T]
	Sma10Roc20   *trend.Sma[T]
	Sma15Roc30   *trend.Sma[T]
	Sma50Roc40   *trend.Sma[T]
	Sma65Roc65   *trend.Sma[T]
	Sma75Roc75   *trend.Sma[T]
	Sma100Roc100 *trend.Sma[T]
	Sma130Roc195 *trend.Sma[T]
	Sma130Roc265 *trend.Sma[T]
	Sma130Roc390 *trend.Sma[T]
	Sma195Roc530 *trend.Sma[T]
}

// NewPringsSpecialK function initializes a new Martin Pring's Special K instance.
func NewPringsSpecialK[T helper.Float]() *PringsSpecialK[T] {
	return &PringsSpecialK[T]{
		Roc10:  trend.NewRocWithPeriod[T](10),
		Roc15:  trend.NewRocWithPeriod[T](15),
		Roc20:  trend.NewRocWithPeriod[T](20),
		Roc30:  trend.NewRocWithPeriod[T](30),
		Roc40:  trend.NewRocWithPeriod[T](40),
		Roc65:  trend.NewRocWithPeriod[T](65),
		Roc75:  trend.NewRocWithPeriod[T](75),
		Roc100: trend.NewRocWithPeriod[T](100),
		Roc195: trend.NewRocWithPeriod[T](195),
		Roc265: trend.NewRocWithPeriod[T](265),
		Roc390: trend.NewRocWithPeriod[T](390),
		Roc530: trend.NewRocWithPeriod[T](530),

		Sma10Roc10:   trend.NewSmaWithPeriod[T](10),
		Sma10Roc15:   trend.NewSmaWithPeriod[T](10),
		Sma10Roc20:   trend.NewSmaWithPeriod[T](10),
		Sma15Roc30:   trend.NewSmaWithPeriod[T](15),
		Sma50Roc40:   trend.NewSmaWithPeriod[T](50),
		Sma65Roc65:   trend.NewSmaWithPeriod[T](65),
		Sma75Roc75:   trend.NewSmaWithPeriod[T](75),
		Sma100Roc100: trend.NewSmaWithPeriod[T](100),
		Sma130Roc195: trend.NewSmaWithPeriod[T](130),
		Sma130Roc265: trend.NewSmaWithPeriod[T](130),
		Sma130Roc390: trend.NewSmaWithPeriod[T](130),
		Sma195Roc530: trend.NewSmaWithPeriod[T](195),
	}
}

func (p *PringsSpecialK[T]) Compute(closings <-chan T) <-chan T {
	c := helper.Duplicate(closings, 12)

	roc10 := p.Roc10.Compute(c[0])
	roc15 := p.Roc15.Compute(c[1])
	roc20 := p.Roc20.Compute(c[2])
	roc30 := p.Roc30.Compute(c[3])
	roc40 := p.Roc40.Compute(c[4])
	roc65 := p.Roc65.Compute(c[5])
	roc75 := p.Roc75.Compute(c[6])
	roc100 := p.Roc100.Compute(c[7])
	roc195 := p.Roc195.Compute(c[8])
	roc265 := p.Roc265.Compute(c[9])
	roc390 := p.Roc390.Compute(c[10])
	roc530 := p.Roc530.Compute(c[11])

	sma10Roc10 := p.Sma10Roc10.Compute(roc10)
	sma10Roc15 := p.Sma10Roc15.Compute(roc15)
	sma10Roc20 := p.Sma10Roc20.Compute(roc20)
	sma15Roc30 := p.Sma15Roc30.Compute(roc30)
	sma50Roc40 := p.Sma50Roc40.Compute(roc40)
	sma65Roc65 := p.Sma65Roc65.Compute(roc65)
	sma75Roc75 := p.Sma75Roc75.Compute(roc75)
	sma100Roc100 := p.Sma100Roc100.Compute(roc100)
	sma130Roc195 := p.Sma130Roc195.Compute(roc195)
	sma130Roc265 := p.Sma130Roc265.Compute(roc265)
	sma130Roc390 := p.Sma130Roc390.Compute(roc390)
	sma195Roc530 := p.Sma195Roc530.Compute(roc530)

	maxIdle := p.Sma195Roc530.IdlePeriod() + p.Roc530.IdlePeriod()

	sma10Roc10 = helper.Skip(sma10Roc10, maxIdle-p.Sma10Roc10.IdlePeriod()-p.Roc10.IdlePeriod())
	sma10Roc15 = helper.Skip(sma10Roc15, maxIdle-p.Sma10Roc15.IdlePeriod()-p.Roc15.IdlePeriod())
	sma10Roc20 = helper.Skip(sma10Roc20, maxIdle-p.Sma10Roc20.IdlePeriod()-p.Roc20.IdlePeriod())
	sma15Roc30 = helper.Skip(sma15Roc30, maxIdle-p.Sma15Roc30.IdlePeriod()-p.Roc30.IdlePeriod())
	sma50Roc40 = helper.Skip(sma50Roc40, maxIdle-p.Sma50Roc40.IdlePeriod()-p.Roc40.IdlePeriod())
	sma65Roc65 = helper.Skip(sma65Roc65, maxIdle-p.Sma65Roc65.IdlePeriod()-p.Roc65.IdlePeriod())
	sma75Roc75 = helper.Skip(sma75Roc75, maxIdle-p.Sma75Roc75.IdlePeriod()-p.Roc75.IdlePeriod())
	sma100Roc100 = helper.Skip(sma100Roc100, maxIdle-p.Sma100Roc100.IdlePeriod()-p.Roc100.IdlePeriod())
	sma130Roc195 = helper.Skip(sma130Roc195, maxIdle-p.Sma130Roc195.IdlePeriod()-p.Roc195.IdlePeriod())
	sma130Roc265 = helper.Skip(sma130Roc265, maxIdle-p.Sma130Roc265.IdlePeriod()-p.Roc265.IdlePeriod())
	sma130Roc390 = helper.Skip(sma130Roc390, maxIdle-p.Sma130Roc390.IdlePeriod()-p.Roc390.IdlePeriod())

	p0 := helper.MultiplyBy(sma10Roc10, 1)
	p1 := helper.Add(p0, helper.MultiplyBy(sma10Roc15, 2))
	p2 := helper.Add(p1, helper.MultiplyBy(sma10Roc20, 3))
	p3 := helper.Add(p2, helper.MultiplyBy(sma15Roc30, 4))
	p4 := helper.Add(p3, helper.MultiplyBy(sma50Roc40, 1))
	p5 := helper.Add(p4, helper.MultiplyBy(sma65Roc65, 2))
	p6 := helper.Add(p5, helper.MultiplyBy(sma75Roc75, 3))
	p7 := helper.Add(p6, helper.MultiplyBy(sma100Roc100, 4))
	p8 := helper.Add(p7, helper.MultiplyBy(sma130Roc195, 1))
	p9 := helper.Add(p8, helper.MultiplyBy(sma130Roc265, 2))
	p10 := helper.Add(p9, helper.MultiplyBy(sma130Roc390, 3))
	p11 := helper.Add(p10, helper.MultiplyBy(sma195Roc530, 4))

	return p11
}
