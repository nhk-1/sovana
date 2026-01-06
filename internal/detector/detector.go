package detector

import (
	"math"
	"sort"
	"strings"

	"sovana/internal/model"
)

// DetectSubscriptions finds recurring transactions that look like subscriptions.
func DetectSubscriptions(transactions []model.Transaction) []model.Subscription {
	normalized := groupByMerchant(transactions)
	subs := make([]model.Subscription, 0)

	for merchant, txs := range normalized {
		if len(txs) < 2 {
			continue
		}
		sort.Slice(txs, func(i, j int) bool { return txs[i].Date.Before(txs[j].Date) })

		if !isMonthly(txs) {
			continue
		}

		avg := averageAmount(txs)
		if !amountsSimilar(txs, avg) {
			continue
		}

		first := txs[0].Date
		last := txs[len(txs)-1].Date
		subs = append(subs, model.Subscription{
			Merchant:   merchant,
			Amount:     round(avg, 2),
			Period:     "monthly",
			FirstSeen:  first,
			LastSeen:   last,
			AnnualCost: round(avg*12, 2),
			CancelURL:  cancelURL(merchant),
		})
	}

	sort.Slice(subs, func(i, j int) bool { return subs[i].Merchant < subs[j].Merchant })
	return subs
}

func groupByMerchant(transactions []model.Transaction) map[string][]model.Transaction {
	groups := make(map[string][]model.Transaction)
	for _, tx := range transactions {
		merchant := normalizeLabel(tx.Label)
		groups[merchant] = append(groups[merchant], tx)
	}
	return groups
}

func normalizeLabel(label string) string {
	lower := strings.ToLower(label)
	replacements := []string{"payment", "carte", "cb", "prlv", "debit", "visa"}
	for _, r := range replacements {
		lower = strings.ReplaceAll(lower, r, "")
	}
	lower = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == ' ' {
			return r
		}
		return -1
	}, lower)
	return strings.TrimSpace(strings.Join(strings.Fields(lower), " "))
}

func isMonthly(txs []model.Transaction) bool {
	if len(txs) < 2 {
		return false
	}
	intervals := make([]int, 0, len(txs)-1)
	for i := 1; i < len(txs); i++ {
		days := int(math.Abs(txs[i].Date.Sub(txs[i-1].Date).Hours() / 24))
		intervals = append(intervals, days)
	}
	sort.Ints(intervals)
	median := intervals[len(intervals)/2]
	return median >= 28 && median <= 31
}

func averageAmount(txs []model.Transaction) float64 {
	var sum float64
	for _, tx := range txs {
		sum += math.Abs(tx.Amount)
	}
	return sum / float64(len(txs))
}

func amountsSimilar(txs []model.Transaction, avg float64) bool {
	tolerance := math.Max(1.0, avg*0.1)
	for _, tx := range txs {
		if math.Abs(math.Abs(tx.Amount)-avg) > tolerance {
			return false
		}
	}
	return true
}

func round(value float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	return math.Round(value*pow) / pow
}

func cancelURL(merchant string) string {
	known := map[string]string{
		"netflix":      "https://www.netflix.com/cancel",
		"spotify":      "https://www.spotify.com/account/subscription",
		"apple":        "https://support.apple.com/HT202039",
		"amazon prime": "https://www.amazon.com/gp/video/settings",
		"adobe":        "https://account.adobe.com/plans",
		"canal":        "https://client.canalplus.com/compte/resiliation",
	}
	key := strings.ToLower(merchant)
	if url, ok := known[key]; ok {
		return url
	}
	return ""
}
