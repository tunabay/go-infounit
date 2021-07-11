// Copyright (c) 2020 Hirotsuna Mizuno. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package infounit

import (
	"math"
	"strconv"
)

//
const (
	kilo uint64 = 1000        // SI prefix k, kilo
	mega        = 1000 * kilo // SI prefix M, mega
	giga        = 1000 * mega // SI prefix G, giga
	tera        = 1000 * giga // SI prefix T, tera
	peta        = 1000 * tera // SI prefix P, peta
	exa         = 1000 * peta // SI prefix E, exa
	kibi        = 1024        // Binary prefix Ki, kibi
	mebi        = 1024 * kibi // Binary prefix Mi, mebi
	gibi        = 1024 * mebi // Binary prefix Gi, gibi
	tebi        = 1024 * gibi // Binary prefix Ti, tebi
	pebi        = 1024 * tebi // Binary prefix Pi, pebi
	exbi        = 1024 * pebi // Binary prefix Ei, exbi
)

//
type prefix struct {
	thresholds [6]uint64
	preAbbr    [6]string
	preFull    [6]string
}

//
var (
	siPrefix = &prefix{
		thresholds: [6]uint64{kilo, mega, giga, tera, peta, exa},
		preAbbr:    [6]string{"k", "M", "G", "T", "P", "E"},
		preFull:    [6]string{"kilo", "mega", "giga", "tera", "peta", "exa"},
	}
	binPrefix = &prefix{
		thresholds: [6]uint64{kibi, mebi, gibi, tebi, pebi, exbi},
		preAbbr:    [6]string{"Ki", "Mi", "Gi", "Ti", "Pi", "Ei"},
		preFull:    [6]string{"kibi", "mebi", "gibi", "tebi", "pebi", "exbi"},
	}
)

// formatUint is used by both ByteCount and BitCount.
func (p *prefix) formatUint(v uint64, precision int, full, space bool, uAbbr, uFull string) string {
	pre, unit, pls := p.preAbbr, uAbbr, ""
	if full {
		pre, unit, pls = p.preFull, uFull, "s"
	}
	sp := ""
	if space {
		sp = " "
	}

	if v == 1 {
		pls = ""
	}
	if v < p.thresholds[0] {
		return strconv.FormatUint(v, 10) + sp + unit + pls
	}
	var ret string
	for i := 0; i < 6; i++ {
		if i < 5 && p.thresholds[i+1] <= v {
			continue
		}
		if v == p.thresholds[i] {
			pls = ""
		}
		bv := float64(v) / float64(p.thresholds[i])
		ret = strconv.FormatFloat(bv, 'f', precision, 64) + sp + pre[i] + unit + pls
		break
	}
	return ret
}

// formatFloat is used by ByteRate.
func (p *prefix) formatFloat(v float64, precision int, full, space bool, uAbbr, sufAbbr string) string {
	pre, unit, pls, suf := p.preAbbr, uAbbr, "", sufAbbr
	if full {
		pre, unit, pls, suf = p.preFull, unitBitRateFull, "s", " "+unitBitRateLongSuffix
	}
	sp := ""
	if space {
		sp = " "
	}

	if v == 1.00 {
		pls = ""
	}
	if math.IsNaN(v) || math.IsInf(v, +1) || math.IsInf(v, -1) {
		return strconv.FormatFloat(v, 'f', precision, 64) + sp + unit + pls + suf
	}
	if v < float64(p.thresholds[0]) {
		return strconv.FormatFloat(v, 'f', precision, 64) + sp + unit + pls + suf
	}
	var ret string
	for i := 0; i < 6; i++ {
		if i < 5 && float64(p.thresholds[i+1]) <= v {
			continue
		}
		if v == float64(p.thresholds[i]) {
			pls = ""
		}
		bv := v / float64(p.thresholds[i])
		ret = strconv.FormatFloat(bv, 'f', precision, 64) + sp + pre[i] + unit + pls + suf
		break
	}
	return ret
}

// jsonNULL is the null expression in JSON.
const jsonNULL = "null"
