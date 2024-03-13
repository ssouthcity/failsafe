package kensoy

type activity struct {
	name  string
	image string
}

var abbreviations = map[string]activity{
	"ron": {"Root of Nightmares", "assets/raids/ron.jpg"},
	"vow": {"Vow of the Disciple", "assets/raids/vow.jpg"},
}

func abbreviationToActivity(abbrv string) (activity, bool) {
	activity, ok := abbreviations[abbrv]
	return activity, ok
}

func findActivityMentioned(message string) (activity, bool) {
	tokens := tokenize(message)

	var (
		act activity
		ok  bool
	)

	for _, token := range tokens {
		if act, ok = abbreviationToActivity(token); ok {
			return act, true
		}
	}

	return act, false
}
