---
title: Prometheus
---

## Prometheus


## functions
* hash: 8fdfa8abeaaaf021d3ca614e4b50af317b1f4145
* https://github.com/prometheus/prometheus/blob/master/promql/functions.go

``` go
// extrapolatedRate is a utility function for rate/increase/delta.
// It calculates the rate (allowing for counter resets if isCounter is true),
// extrapolates if the first/last sample is close to the boundary, and returns
// the result as either per-second (if isRate is true) or overall.
func extrapolatedRate(vals []Value, args Expressions, enh *EvalNodeHelper, isCounter bool, isRate bool) Vector {
	ms := args[0].(*MatrixSelector)

	var (
		matrix     = vals[0].(Matrix)
		rangeStart = enh.ts - durationMilliseconds(ms.Range+ms.Offset)
		rangeEnd   = enh.ts - durationMilliseconds(ms.Offset)
	)

	for _, samples := range matrix {
		// No sense in trying to compute a rate without at least two points. Drop
		// this Vector element.
		if len(samples.Points) < 2 {
			continue
		}
		var (
			counterCorrection float64
			lastValue         float64
		)
		for _, sample := range samples.Points {
		    // comment: counter で0 に戻ったときの対策として、counterCorrection をおいているのだと思われ 
			if isCounter && sample.V < lastValue {
				counterCorrection += lastValue
			}
			lastValue = sample.V
		}
		resultValue := lastValue - samples.Points[0].V + counterCorrection

		// Duration between first/last samples and boundary of range.
		durationToStart := float64(samples.Points[0].T-rangeStart) / 1000
		durationToEnd := float64(rangeEnd-samples.Points[len(samples.Points)-1].T) / 1000

		sampledInterval := float64(samples.Points[len(samples.Points)-1].T-samples.Points[0].T) / 1000
		averageDurationBetweenSamples := sampledInterval / float64(len(samples.Points)-1)

		if isCounter && resultValue > 0 && samples.Points[0].V >= 0 {
			// Counters cannot be negative. If we have any slope at
			// all (i.e. resultValue went up), we can extrapolate
			// the zero point of the counter. If the duration to the
			// zero point is shorter than the durationToStart, we
			// take the zero point as the start of the series,
			// thereby avoiding extrapolation to negative counter
			// values.
			// comment:  
			durationToZero := sampledInterval * (samples.Points[0].V / resultValue)
			if durationToZero < durationToStart {
				durationToStart = durationToZero
			}
		}

		// If the first/last samples are close to the boundaries of the range,
		// extrapolate the result. This is as we expect that another sample
		// will exist given the spacing between samples we've seen thus far,
		// with an allowance for noise.
		extrapolationThreshold := averageDurationBetweenSamples * 1.1
		extrapolateToInterval := sampledInterval

		if durationToStart < extrapolationThreshold {
			extrapolateToInterval += durationToStart
		} else {
			extrapolateToInterval += averageDurationBetweenSamples / 2
		}
		if durationToEnd < extrapolationThreshold {
			extrapolateToInterval += durationToEnd
		} else {
			extrapolateToInterval += averageDurationBetweenSamples / 2
		}
		resultValue = resultValue * (extrapolateToInterval / sampledInterval)
		if isRate {
			resultValue = resultValue / ms.Range.Seconds()
		}

		enh.out = append(enh.out, Sample{
			Point: Point{V: resultValue},
		})
	}
	return enh.out
}

// === rate(node ValueTypeMatrix) Vector ===
func funcRate(vals []Value, args Expressions, enh *EvalNodeHelper) Vector {
	return extrapolatedRate(vals, args, enh, true, true)
}
```

rate と increase は、
時間範囲の開始から最初のサンプルまでの間隔と最後のサンプルから時間範囲 の終了までの間隔が
データの平均インターバルの 110% 以内であれば、時系列データが時間範囲の境界 を越えて続く前提で動作する。
そうでない場合、手持ちのサンプルを超えて平均インターバルの 50% に時系列データがあり、値がゼロ未満になることはない前提で動作する。

## Reference
* [instrumenting-go-applications](https://banzaicloud.com/blog/instrumenting-go-applications/)
* [EXPLORING PROMETHEUS GO CLIENT METRICS](https://povilasv.me/prometheus-go-metrics/)
  * メトリクスの解説
