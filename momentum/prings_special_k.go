package momentum

import (
	"github.com/cinar/indicator/v2/helper"
	"github.com/cinar/indicator/v2/trend"
)

type PringsSpecialK[T helper.Float] struct {
	Roc10  *trend.Roc[T]
	Roc15  *trend.Roc[T]
	Roc50  *trend.Roc[T]
	Roc65  *trend.Roc[T]
	Roc75  *trend.Roc[T]
	Roc100 *trend.Roc[T]
	Roc130 *trend.Roc[T]
	Roc195 *trend.Roc[T]

	Sma10  *trend.Sma[T]
	Sma15  *trend.Sma[T]
	Sma20  *trend.Sma[T]
	Sma30  *trend.Sma[T]
	Sma40  *trend.Sma[T]
	Sma65  *trend.Sma[T]
	Sma75  *trend.Sma[T]
	Sma100 *trend.Sma[T]
	Sma195 *trend.Sma[T]
	Sma265 *trend.Sma[T]
	Sma390 *trend.Sma[T]
	Sma530 *trend.Sma[T]
}

// NewPringsSpecialK function initializes a new Martin Pring's Special K instance.
func NewPringsSpecialK[T helper.Float]() *PringsSpecialK[T] {
	return &PringsSpecialK[T]{
		Roc10:  trend.NewRocWithPeriod[T](10),
		Roc15:  trend.NewRocWithPeriod[T](15),
		Roc50:  trend.NewRocWithPeriod[T](50),
		Roc65:  trend.NewRocWithPeriod[T](65),
		Roc75:  trend.NewRocWithPeriod[T](75),
		Roc100: trend.NewRocWithPeriod[T](100),
		Roc130: trend.NewRocWithPeriod[T](130),
		Roc195: trend.NewRocWithPeriod[T](195),

		Sma10:  trend.NewSmaWithPeriod[T](10),
		Sma15:  trend.NewSmaWithPeriod[T](15),
		Sma20:  trend.NewSmaWithPeriod[T](20),
		Sma30:  trend.NewSmaWithPeriod[T](30),
		Sma40:  trend.NewSmaWithPeriod[T](40),
		Sma65:  trend.NewSmaWithPeriod[T](65),
		Sma75:  trend.NewSmaWithPeriod[T](75),
		Sma100: trend.NewSmaWithPeriod[T](100),
		Sma195: trend.NewSmaWithPeriod[T](195),
		Sma265: trend.NewSmaWithPeriod[T](265),
		Sma390: trend.NewSmaWithPeriod[T](390),
		Sma530: trend.NewSmaWithPeriod[T](530),
	}
}

func (p *PringsSpecialK[T]) Compute(closings <-chan T) <-chan T {
	c := helper.Duplicate(closings, 8)

	roc10 := p.Roc10.Compute(c[0])
	roc15 := p.Roc15.Compute(c[1])
	roc50 := p.Roc50.Compute(c[2])
	roc65 := p.Roc65.Compute(c[3])
	roc75 := p.Roc75.Compute(c[4])
	roc100 := p.Roc100.Compute(c[5])
	roc130 := p.Roc130.Compute(c[6])
	roc195 := p.Roc195.Compute(c[7])

	roc10s := helper.Duplicate(roc10, 3)
	sma10 := p.Sma10.Compute(roc10s[0])
	sma15 := p.Sma15.Compute(roc10s[1])
	sma20 := p.Sma20.Compute(roc10s[2])
	sma30 := p.Sma30.Compute(roc15)
	sma40 := p.Sma40.Compute(roc50)
	sma65 := p.Sma65.Compute(roc65)
	sma75 := p.Sma75.Compute(roc75)
	sma100 := p.Sma100.Compute(roc100)
	roc130s := helper.Duplicate(roc130, 3)
	sma195 := p.Sma195.Compute(roc130s[0])
	sma265 := p.Sma265.Compute(roc130s[1])
	sma390 := p.Sma390.Compute(roc130s[2])
	sma530 := p.Sma530.Compute(roc195)

	maxIdle := p.Sma530.IdlePeriod() + p.Roc195.IdlePeriod()

	sma10 = helper.Skip(sma10, maxIdle-p.Sma10.IdlePeriod()-p.Roc10.IdlePeriod())
	sma15 = helper.Skip(sma15, maxIdle-p.Sma15.IdlePeriod()-p.Roc10.IdlePeriod())
	sma20 = helper.Skip(sma20, maxIdle-p.Sma20.IdlePeriod()-p.Roc10.IdlePeriod())
	sma30 = helper.Skip(sma30, maxIdle-p.Sma30.IdlePeriod()-p.Roc15.IdlePeriod())
	sma40 = helper.Skip(sma40, maxIdle-p.Sma40.IdlePeriod()-p.Roc50.IdlePeriod())
	sma65 = helper.Skip(sma65, maxIdle-p.Sma65.IdlePeriod()-p.Roc65.IdlePeriod())
	sma75 = helper.Skip(sma75, maxIdle-p.Sma75.IdlePeriod()-p.Roc75.IdlePeriod())
	sma100 = helper.Skip(sma100, maxIdle-p.Sma100.IdlePeriod()-p.Roc100.IdlePeriod())
	sma195 = helper.Skip(sma195, maxIdle-p.Sma195.IdlePeriod()-p.Roc130.IdlePeriod())
	sma265 = helper.Skip(sma265, maxIdle-p.Sma265.IdlePeriod()-p.Roc130.IdlePeriod())
	sma390 = helper.Skip(sma390, maxIdle-p.Sma390.IdlePeriod()-p.Roc130.IdlePeriod())
	//sma530 = helper.Skip(sma530, p.Sma530.IdlePeriod()-p.Sma530.IdlePeriod())

	p0 := helper.MultiplyBy(sma10, 1)
	p1 := helper.Add(p0, helper.MultiplyBy(sma15, 2))
	p2 := helper.Add(p1, helper.MultiplyBy(sma20, 3))
	p3 := helper.Add(p2, helper.MultiplyBy(sma30, 4))
	p4 := helper.Add(p3, helper.MultiplyBy(sma40, 1))
	p5 := helper.Add(p4, helper.MultiplyBy(sma65, 2))
	p6 := helper.Add(p5, helper.MultiplyBy(sma75, 3))
	p7 := helper.Add(p6, helper.MultiplyBy(sma100, 4))
	p8 := helper.Add(p7, helper.MultiplyBy(sma195, 1))
	p9 := helper.Add(p8, helper.MultiplyBy(sma265, 2))
	p10 := helper.Add(p9, helper.MultiplyBy(sma390, 3))
	p11 := helper.Add(p10, helper.MultiplyBy(sma530, 4))

	return p11
}
