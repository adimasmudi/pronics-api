package helper

type GoogleUser struct {
	Id            string
	Email         string
	VerifiedEmail bool
	Picture       string
}

type DistanceMatrixResult struct {
	DestinationAddress []string          `json:"destination_addresses"`
	OriginAddress      []string          `json:"origin_addresses"`
	Rows               []ElementDistance `json:"rows"`
}

type ElementDistance struct {
	Elements []DistanceAndDuration `json:"elements"`
}

type DistanceAndDuration struct {
	Distance DistanceElement `json:"distance"`
	Duration DurationElement `json:"duration"`
}

type DistanceElement struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}

type DurationElement struct {
	Text  string `json:"text"`
	Value int    `json:"value"`
}
