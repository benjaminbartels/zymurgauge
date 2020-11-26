//nolint:maligned
package brewfather

type Recipe struct {
	AvgWeightedHopstandTemp  int              `json:"avgWeightedHopstandTemp"`
	ID                       string           `json:"_id"`
	BoilSize                 float64          `json:"boilSize"`
	Og                       float64          `json:"og"`
	Tags                     interface{}      `json:"tags"`
	PreBoilGravity           float64          `json:"preBoilGravity"`
	Path                     string           `json:"path"`
	SearchTags               []interface{}    `json:"searchTags"`
	TimestampMs              int64            `json:"_timestamp_ms"`
	HopsTotalAmount          float64          `json:"hopsTotalAmount"`
	StyleAbv                 bool             `json:"styleAbv"`
	Attenuation              float64          `json:"attenuation"`
	Mash                     Mash             `json:"mash"`
	Abv                      float64          `json:"abv"`
	Fermentation             Fermentation     `json:"fermentation"`
	PrimaryTemp              float64          `json:"primaryTemp"`
	Color                    float64          `json:"color"`
	StyleIbu                 bool             `json:"styleIbu"`
	Yeasts                   []Yeasts         `json:"yeasts"`
	SumDryHopPerLiter        int              `json:"sumDryHopPerLiter"`
	StyleRbr                 bool             `json:"styleRbr"`
	BoilTime                 int              `json:"boilTime"`
	Water                    Water            `json:"water"`
	Defaults                 Defaults         `json:"defaults"`
	Ibu                      float64          `json:"ibu"`
	FgEstimated              float64          `json:"fgEstimated"`
	FermentablesTotalAmount  float64          `json:"fermentablesTotalAmount"`
	StyleColor               bool             `json:"styleColor"`
	JSONType                 string           `json:"_type"`
	Created                  Created          `json:"_created"`
	StyleBuGu                bool             `json:"styleBuGu"`
	Hops                     []Hops           `json:"hops"`
	Timestamp                Timestamp        `json:"_timestamp"`
	Fermentables             []Fermentables   `json:"fermentables"`
	HopStandMinutes          int              `json:"hopStandMinutes"`
	FgFormula                string           `json:"fgFormula"`
	BatchSize                float64          `json:"batchSize"`
	StyleFg                  bool             `json:"styleFg"`
	StyleOg                  bool             `json:"styleOg"`
	Init                     bool             `json:"_init"`
	Efficiency               float64          `json:"efficiency"`
	Style                    Style            `json:"style"`
	ManualFg                 bool             `json:"manualFg"`
	CarbonationStyle         CarbonationStyle `json:"carbonationStyle"`
	Fg                       float64          `json:"fg"`
	Name                     string           `json:"name"`
	Data                     Data             `json:"data"`
	StyleConformity          bool             `json:"styleConformity"`
	FirstWortGravity         interface{}      `json:"firstWortGravity"`
	Author                   interface{}      `json:"author"`
	PostBoilGravity          float64          `json:"postBoilGravity"`
	RbRatio                  float64          `json:"rbRatio"`
	Hidden                   bool             `json:"hidden"`
	BuGuRatio                float64          `json:"buGuRatio"`
	IbuFormula               string           `json:"ibuFormula"`
	FermentableIbu           int              `json:"fermentableIbu"`
	Carbonation              float64          `json:"carbonation"`
	MashEfficiency           float64          `json:"mashEfficiency"`
	Miscs                    []Miscs          `json:"miscs"`
	Equipment                Equipment        `json:"equipment"`
	Rev                      string           `json:"_rev"`
	Nutrition                Nutrition        `json:"nutrition"`
	Version                  string           `json:"_version"`
	StyleCarb                bool             `json:"styleCarb"`
	YeastToleranceExceededBy interface{}      `json:"yeastToleranceExceededBy"`
	Type                     string           `json:"type"`
	OgPlato                  float64          `json:"ogPlato"`
}

type MashStep struct {
	DisplayStepTemp float64     `json:"displayStepTemp"`
	RampTime        interface{} `json:"rampTime"`
	StepTemp        float64     `json:"stepTemp"`
	StepTime        int         `json:"stepTime"`
	Type            string      `json:"type"`
	Name            string      `json:"name"`
}

type Mash struct {
	Steps []MashStep `json:"steps"`
	Name  string     `json:"name"`
	ID    string     `json:"_id"`
}

type FermentationStep struct {
	Pressure        interface{} `json:"pressure"`
	Ramp            interface{} `json:"ramp"`
	Type            string      `json:"type"`
	StepTemp        float64     `json:"stepTemp"`
	StepTime        int         `json:"stepTime"`
	DisplayPressure interface{} `json:"displayPressure"`
	DisplayStepTemp float64     `json:"displayStepTemp"`
}

type Fermentation struct {
	Steps []FermentationStep `json:"steps"`
	ID    string             `json:"_id"`
	Name  string             `json:"name"`
}

type Yeasts struct {
	Unit              string      `json:"unit"`
	ID                string      `json:"_id"`
	BestBeforeDate    interface{} `json:"bestBeforeDate"`
	Flocculation      string      `json:"flocculation"`
	UserNotes         string      `json:"userNotes"`
	MinAttenuation    interface{} `json:"minAttenuation"`
	Name              string      `json:"name"`
	ProductID         string      `json:"productId"`
	Type              string      `json:"type"`
	MaxTemp           float64     `json:"maxTemp"`
	ManufacturingDate interface{} `json:"manufacturingDate"`
	Laboratory        string      `json:"laboratory"`
	FermentsAll       bool        `json:"fermentsAll"`
	Amount            int         `json:"amount"`
	Attenuation       int         `json:"attenuation"`
	MaxAttenuation    interface{} `json:"maxAttenuation"`
	Form              string      `json:"form"`
	MinTemp           float64     `json:"minTemp"`
	MaxAbv            interface{} `json:"maxAbv"`
	Description       string      `json:"description"`
}

type Created struct {
	Seconds     int `json:"_seconds"`
	Nanoseconds int `json:"_nanoseconds"`
}

type Timestamp struct {
	Seconds     int `json:"_seconds"`
	Nanoseconds int `json:"_nanoseconds"`
}

type Total struct {
	BicarbonateMeqL            float64   `json:"bicarbonateMeqL"`
	Rev                        string    `json:"_rev"`
	Anions                     float64   `json:"anions"`
	Chloride                   float64   `json:"chloride"`
	Hardness                   int       `json:"hardness"`
	Hidden                     bool      `json:"hidden"`
	Sodium                     int       `json:"sodium"`
	Name                       string    `json:"name"`
	Bicarbonate                int       `json:"bicarbonate"`
	Ph                         float64   `json:"ph"`
	Type                       string    `json:"type"`
	Calcium                    float64   `json:"calcium"`
	Created                    Created   `json:"_created"`
	Version                    string    `json:"_version"`
	Alkalinity                 float64   `json:"alkalinity"`
	ResidualAlkalinity         float64   `json:"residualAlkalinity"`
	IonBalance                 int       `json:"ionBalance"`
	IonBalanceOff              bool      `json:"ionBalanceOff"`
	Cations                    float64   `json:"cations"`
	Timestamp                  Timestamp `json:"_timestamp"`
	ID                         string    `json:"_id"`
	ResidualAlkalinityMeqLCalc float64   `json:"residualAlkalinityMeqLCalc"`
	TimestampMs                int64     `json:"_timestamp_ms"`
	SoClRatio                  int       `json:"soClRatio"`
	Sulfate                    float64   `json:"sulfate"`
	Magnesium                  float64   `json:"magnesium"`
}

type WaterMash struct {
	Type                       string    `json:"type"`
	Bicarbonate                int       `json:"bicarbonate"`
	Version                    string    `json:"_version"`
	Chloride                   float64   `json:"chloride"`
	ResidualAlkalinityMeqLCalc float64   `json:"residualAlkalinityMeqLCalc"`
	TimestampMs                int64     `json:"_timestamp_ms"`
	Sodium                     int       `json:"sodium"`
	ResidualAlkalinity         float64   `json:"residualAlkalinity"`
	Sulfate                    float64   `json:"sulfate"`
	Created                    Created   `json:"_created"`
	SoClRatio                  int       `json:"soClRatio"`
	Rev                        string    `json:"_rev"`
	BicarbonateMeqL            float64   `json:"bicarbonateMeqL"`
	Cations                    float64   `json:"cations"`
	Ph                         float64   `json:"ph"`
	ID                         string    `json:"_id"`
	Hidden                     bool      `json:"hidden"`
	Alkalinity                 float64   `json:"alkalinity"`
	Magnesium                  float64   `json:"magnesium"`
	IonBalanceOff              bool      `json:"ionBalanceOff"`
	Anions                     float64   `json:"anions"`
	Hardness                   int       `json:"hardness"`
	Timestamp                  Timestamp `json:"_timestamp"`
	Calcium                    float64   `json:"calcium"`
	IonBalance                 int       `json:"ionBalance"`
	Name                       string    `json:"name"`
}

type SourceTargetDiff struct {
	SoClRatio                  float64 `json:"soClRatio"`
	Sodium                     int     `json:"sodium"`
	Anions                     float64 `json:"anions"`
	Hardness                   int     `json:"hardness"`
	Bicarbonate                int     `json:"bicarbonate"`
	Sulfate                    int     `json:"sulfate"`
	ResidualAlkalinityMeqLCalc float64 `json:"residualAlkalinityMeqLCalc"`
	IonBalance                 int     `json:"ionBalance"`
	ResidualAlkalinity         float64 `json:"residualAlkalinity"`
	Alkalinity                 float64 `json:"alkalinity"`
	Chloride                   int     `json:"chloride"`
	Calcium                    int     `json:"calcium"`
	Cations                    float64 `json:"cations"`
	BicarbonateMeqL            float64 `json:"bicarbonateMeqL"`
	IonBalanceOff              bool    `json:"ionBalanceOff"`
	Magnesium                  int     `json:"magnesium"`
}

type SpargeTargetDiff struct {
	Calcium                    float64 `json:"calcium"`
	Cations                    float64 `json:"cations"`
	SoClRatio                  float64 `json:"soClRatio"`
	ResidualAlkalinityMeqLCalc float64 `json:"residualAlkalinityMeqLCalc"`
	ResidualAlkalinity         float64 `json:"residualAlkalinity"`
	Magnesium                  float64 `json:"magnesium"`
	Hardness                   int     `json:"hardness"`
	Sodium                     int     `json:"sodium"`
	IonBalanceOff              bool    `json:"ionBalanceOff"`
	IonBalance                 int     `json:"ionBalance"`
	Chloride                   float64 `json:"chloride"`
	BicarbonateMeqL            float64 `json:"bicarbonateMeqL"`
	Bicarbonate                int     `json:"bicarbonate"`
	Sulfate                    float64 `json:"sulfate"`
	Anions                     float64 `json:"anions"`
	Alkalinity                 float64 `json:"alkalinity"`
}

type TotalTargetDiff struct {
	Cations                    float64 `json:"cations"`
	Anions                     float64 `json:"anions"`
	ResidualAlkalinityMeqLCalc float64 `json:"residualAlkalinityMeqLCalc"`
	IonBalanceOff              bool    `json:"ionBalanceOff"`
	Hardness                   int     `json:"hardness"`
	Sulfate                    float64 `json:"sulfate"`
	SoClRatio                  float64 `json:"soClRatio"`
	Magnesium                  float64 `json:"magnesium"`
	Calcium                    float64 `json:"calcium"`
	ResidualAlkalinity         float64 `json:"residualAlkalinity"`
	Bicarbonate                int     `json:"bicarbonate"`
	Chloride                   float64 `json:"chloride"`
	Alkalinity                 float64 `json:"alkalinity"`
	IonBalance                 int     `json:"ionBalance"`
	Sodium                     int     `json:"sodium"`
	BicarbonateMeqL            float64 `json:"bicarbonateMeqL"`
}

type MashAcid struct {
	Type           string `json:"type"`
	Concentration  int    `json:"concentration"`
	Amount         int    `json:"amount"`
	AlkalinityMeqL int    `json:"alkalinityMeqL"`
}

type MashAdjustments struct {
	SodiumMetabisulfite    int        `json:"sodiumMetabisulfite"`
	Magnesium              float64    `json:"magnesium"`
	Calcium                float64    `json:"calcium"`
	MagnesiumChloride      int        `json:"magnesiumChloride"`
	CalciumCarbonate       int        `json:"calciumCarbonate"`
	Bicarbonate            int        `json:"bicarbonate"`
	Sulfate                float64    `json:"sulfate"`
	SodiumChloride         int        `json:"sodiumChloride"`
	CalciumHydroxide       int        `json:"calciumHydroxide"`
	SodiumBicarbonate      int        `json:"sodiumBicarbonate"`
	Acids                  []MashAcid `json:"acids"`
	SodiumMetabisulfitePPM int        `json:"sodiumMetabisulfitePPM"`
	CalciumSulfate         float64    `json:"calciumSulfate"`
	Volume                 float64    `json:"volume"`
	CalciumChloride        float64    `json:"calciumChloride"`
	MagnesiumSulfate       float64    `json:"magnesiumSulfate"`
	Sodium                 int        `json:"sodium"`
	Chloride               float64    `json:"chloride"`
}

type TotalAdjustments struct {
	Magnesium              float64 `json:"magnesium"`
	MagnesiumChloride      int     `json:"magnesiumChloride"`
	Chloride               float64 `json:"chloride"`
	CalciumChloride        float64 `json:"calciumChloride"`
	CalciumCarbonate       int     `json:"calciumCarbonate"`
	SodiumChloride         int     `json:"sodiumChloride"`
	SodiumBicarbonate      int     `json:"sodiumBicarbonate"`
	Sulfate                float64 `json:"sulfate"`
	Volume                 float64 `json:"volume"`
	CalciumHydroxide       int     `json:"calciumHydroxide"`
	Calcium                float64 `json:"calcium"`
	CalciumSulfate         float64 `json:"calciumSulfate"`
	SodiumMetabisulfitePPM int     `json:"sodiumMetabisulfitePPM"`
	SodiumMetabisulfite    int     `json:"sodiumMetabisulfite"`
	Bicarbonate            int     `json:"bicarbonate"`
	MagnesiumSulfate       float64 `json:"magnesiumSulfate"`
	Sodium                 int     `json:"sodium"`
}

type Source struct {
	Type                       string    `json:"type"`
	Sodium                     int       `json:"sodium"`
	Rev                        string    `json:"_rev"`
	Anions                     float64   `json:"anions"`
	Version                    string    `json:"_version"`
	TimestampMs                int64     `json:"_timestamp_ms"`
	IonBalance                 int       `json:"ionBalance"`
	ResidualAlkalinity         float64   `json:"residualAlkalinity"`
	Chloride                   int       `json:"chloride"`
	Alkalinity                 float64   `json:"alkalinity"`
	Bicarbonate                int       `json:"bicarbonate"`
	Created                    Created   `json:"_created"`
	ResidualAlkalinityMeqLCalc float64   `json:"residualAlkalinityMeqLCalc"`
	Calcium                    int       `json:"calcium"`
	Magnesium                  int       `json:"magnesium"`
	Hardness                   int       `json:"hardness"`
	BicarbonateMeqL            float64   `json:"bicarbonateMeqL"`
	Cations                    float64   `json:"cations"`
	Name                       string    `json:"name"`
	Timestamp                  Timestamp `json:"_timestamp"`
	IonBalanceOff              bool      `json:"ionBalanceOff"`
	ID                         string    `json:"_id"`
	Hidden                     bool      `json:"hidden"`
	SoClRatio                  float64   `json:"soClRatio"`
	Ph                         float64   `json:"ph"`
	Sulfate                    int       `json:"sulfate"`
}

type SodiumBicarbonate struct {
	Auto   bool `json:"auto"`
	Mash   bool `json:"mash"`
	Sparge bool `json:"sparge"`
}

type CalciumChloride struct {
	Form   string `json:"form"`
	Auto   bool   `json:"auto"`
	Mash   bool   `json:"mash"`
	Sparge bool   `json:"sparge"`
}

type CalciumSulfate struct {
	Mash   bool `json:"mash"`
	Sparge bool `json:"sparge"`
	Auto   bool `json:"auto"`
}

type MagnesiumSulfate struct {
	Sparge bool `json:"sparge"`
	Auto   bool `json:"auto"`
	Mash   bool `json:"mash"`
}

type CalciumHydroxide struct {
	Sparge bool `json:"sparge"`
	Auto   bool `json:"auto"`
	Mash   bool `json:"mash"`
}

type Settings struct {
	SodiumBicarbonate SodiumBicarbonate `json:"sodiumBicarbonate"`
	CalciumChloride   CalciumChloride   `json:"calciumChloride"`
	CalciumSulfate    CalciumSulfate    `json:"calciumSulfate"`
	MagnesiumSulfate  MagnesiumSulfate  `json:"magnesiumSulfate"`
	CalciumHydroxide  CalciumHydroxide  `json:"calciumHydroxide"`
	AdjustSparge      bool              `json:"adjustSparge"`
}

type SpargeAcid struct {
	Concentration int    `json:"concentration"`
	Amount        int    `json:"amount"`
	Type          string `json:"type"`
}

type SpargeAdjustments struct {
	Sulfate                float64      `json:"sulfate"`
	MagnesiumSulfate       float64      `json:"magnesiumSulfate"`
	Volume                 float64      `json:"volume"`
	Bicarbonate            int          `json:"bicarbonate"`
	CalciumChloride        float64      `json:"calciumChloride"`
	Acids                  []SpargeAcid `json:"acids"`
	SodiumBicarbonate      int          `json:"sodiumBicarbonate"`
	CalciumHydroxide       int          `json:"calciumHydroxide"`
	Magnesium              float64      `json:"magnesium"`
	SodiumMetabisulfitePPM int          `json:"sodiumMetabisulfitePPM"`
	Chloride               float64      `json:"chloride"`
	MagnesiumChloride      int          `json:"magnesiumChloride"`
	CalciumSulfate         float64      `json:"calciumSulfate"`
	SodiumMetabisulfite    int          `json:"sodiumMetabisulfite"`
	Calcium                float64      `json:"calcium"`
	CalciumCarbonate       int          `json:"calciumCarbonate"`
	SodiumChloride         int          `json:"sodiumChloride"`
	Sodium                 int          `json:"sodium"`
}

type Target struct {
	Magnesium                  int     `json:"magnesium"`
	Sodium                     int     `json:"sodium"`
	SoClRatio                  float64 `json:"soClRatio"`
	ResidualAlkalinity         float64 `json:"residualAlkalinity"`
	Hardness                   int     `json:"hardness"`
	Name                       string  `json:"name"`
	Type                       string  `json:"type"`
	Alkalinity                 float64 `json:"alkalinity"`
	Calcium                    int     `json:"calcium"`
	BicarbonateMeqL            float64 `json:"bicarbonateMeqL"`
	IonBalanceOff              bool    `json:"ionBalanceOff"`
	Cations                    float64 `json:"cations"`
	ID                         string  `json:"_id"`
	IonBalance                 int     `json:"ionBalance"`
	Sulfate                    int     `json:"sulfate"`
	ResidualAlkalinityMeqLCalc float64 `json:"residualAlkalinityMeqLCalc"`
	Bicarbonate                int     `json:"bicarbonate"`
	Chloride                   int     `json:"chloride"`
	Anions                     float64 `json:"anions"`
}

type Sparge struct {
	Name                       string    `json:"name"`
	ResidualAlkalinity         float64   `json:"residualAlkalinity"`
	Calcium                    float64   `json:"calcium"`
	Created                    Created   `json:"_created"`
	SoClRatio                  int       `json:"soClRatio"`
	Sulfate                    float64   `json:"sulfate"`
	Type                       string    `json:"type"`
	Cations                    float64   `json:"cations"`
	TimestampMs                int64     `json:"_timestamp_ms"`
	Ph                         float64   `json:"ph"`
	Hidden                     bool      `json:"hidden"`
	ID                         string    `json:"_id"`
	BicarbonateMeqL            float64   `json:"bicarbonateMeqL"`
	Rev                        string    `json:"_rev"`
	Version                    string    `json:"_version"`
	Hardness                   int       `json:"hardness"`
	Chloride                   float64   `json:"chloride"`
	Timestamp                  Timestamp `json:"_timestamp"`
	Anions                     float64   `json:"anions"`
	Bicarbonate                int       `json:"bicarbonate"`
	Magnesium                  float64   `json:"magnesium"`
	IonBalance                 int       `json:"ionBalance"`
	IonBalanceOff              bool      `json:"ionBalanceOff"`
	ResidualAlkalinityMeqLCalc float64   `json:"residualAlkalinityMeqLCalc"`
	Sodium                     int       `json:"sodium"`
	Alkalinity                 float64   `json:"alkalinity"`
}

type MashTargetDiff struct {
	Magnesium                  float64 `json:"magnesium"`
	Anions                     float64 `json:"anions"`
	Hardness                   int     `json:"hardness"`
	Calcium                    float64 `json:"calcium"`
	Sulfate                    float64 `json:"sulfate"`
	ResidualAlkalinity         float64 `json:"residualAlkalinity"`
	Bicarbonate                int     `json:"bicarbonate"`
	SoClRatio                  float64 `json:"soClRatio"`
	Alkalinity                 float64 `json:"alkalinity"`
	Cations                    float64 `json:"cations"`
	BicarbonateMeqL            float64 `json:"bicarbonateMeqL"`
	ResidualAlkalinityMeqLCalc float64 `json:"residualAlkalinityMeqLCalc"`
	Chloride                   float64 `json:"chloride"`
	Sodium                     int     `json:"sodium"`
	IonBalanceOff              bool    `json:"ionBalanceOff"`
	IonBalance                 int     `json:"ionBalance"`
}

type Water struct {
	SpargeAcidPhAdjustment  int               `json:"spargeAcidPhAdjustment"`
	Total                   Total             `json:"total"`
	MashPhDistilled         float64           `json:"mashPhDistilled"`
	Mash                    WaterMash         `json:"mash"`
	AcidPhAdjustment        int               `json:"acidPhAdjustment"`
	SourceTargetDiff        SourceTargetDiff  `json:"sourceTargetDiff"`
	SpargeWaterAmount       interface{}       `json:"spargeWaterAmount"`
	SpargeTargetDiff        SpargeTargetDiff  `json:"spargeTargetDiff"`
	EnableSpargeAdjustments bool              `json:"enableSpargeAdjustments"`
	Style                   string            `json:"style"`
	TotalTargetDiff         TotalTargetDiff   `json:"totalTargetDiff"`
	MashAdjustments         MashAdjustments   `json:"mashAdjustments"`
	TotalAdjustments        TotalAdjustments  `json:"totalAdjustments"`
	Source                  Source            `json:"source"`
	Settings                Settings          `json:"settings"`
	SpargeAdjustments       SpargeAdjustments `json:"spargeAdjustments"`
	Target                  Target            `json:"target"`
	EnableAcidAdjustments   bool              `json:"enableAcidAdjustments"`
	Sparge                  Sparge            `json:"sparge"`
	MashWaterAmount         interface{}       `json:"mashWaterAmount"`
	MashPh                  float64           `json:"mashPh"`
	MashTargetDiff          MashTargetDiff    `json:"mashTargetDiff"`
}

type Defaults struct {
	Abv         string `json:"abv"`
	Gravity     string `json:"gravity"`
	GrainColor  string `json:"grainColor"`
	Attenuation string `json:"attenuation"`
	Temp        string `json:"temp"`
	Volume      string `json:"volume"`
	Preferred   string `json:"preferred"`
	Altitude    string `json:"altitude"`
	Pressure    string `json:"pressure"`
	Color       string `json:"color"`
	Weight      string `json:"weight"`
	Carbonation string `json:"carbonation"`
	Ibu         string `json:"ibu"`
	Hop         string `json:"hop"`
}

type Hops struct {
	Usage             string      `json:"usage"`
	Ibu               float64     `json:"ibu"`
	Temp              interface{} `json:"temp"`
	UserNotes         string      `json:"userNotes"`
	Time              int         `json:"time"`
	Notes             string      `json:"notes"`
	ID                string      `json:"_id"`
	Year              interface{} `json:"year"`
	Amount            float64     `json:"amount"`
	UsedIn            string      `json:"usedIn"`
	BestBeforeDate    interface{} `json:"bestBeforeDate"`
	Origin            string      `json:"origin"`
	Substitutes       string      `json:"substitutes"`
	ManufacturingDate interface{} `json:"manufacturingDate"`
	Type              string      `json:"type"`
	Alpha             int         `json:"alpha"`
	Name              string      `json:"name"`
	Use               string      `json:"use"`
	EditFlag          bool        `json:"_editFlag,omitempty"`
}

type Fermentables struct {
	UsedIn              string      `json:"usedIn"`
	ManufacturingDate   interface{} `json:"manufacturingDate"`
	Moisture            interface{} `json:"moisture"`
	Type                string      `json:"type"`
	Origin              string      `json:"origin"`
	PotentialPercentage float64     `json:"potentialPercentage"`
	Time                int         `json:"time,omitempty"`
	Lovibond            float64     `json:"lovibond"`
	Protein             interface{} `json:"protein"`
	Supplier            string      `json:"supplier"`
	Inventory           interface{} `json:"inventory"`
	Notes               string      `json:"notes"`
	CostPerAmount       interface{} `json:"costPerAmount"`
	NotFermentable      bool        `json:"notFermentable"`
	BestBeforeDate      interface{} `json:"bestBeforeDate"`
	IbuPerAmount        interface{} `json:"ibuPerAmount"`
	Amount              float64     `json:"amount"`
	Substitutes         string      `json:"substitutes"`
	ID                  string      `json:"_id"`
	Name                string      `json:"name"`
	Use                 string      `json:"use,omitempty"`
	Attenuation         float64     `json:"attenuation"`
	DiastaticPower      interface{} `json:"diastaticPower"`
	GrainCategory       string      `json:"grainCategory"`
	Percentage          float64     `json:"percentage"`
	Color               int         `json:"color"`
	UserNotes           string      `json:"userNotes"`
	Potential           float64     `json:"potential"`
}

type Style struct {
	LovibondMin      int     `json:"lovibondMin"`
	CarbonationStyle string  `json:"carbonationStyle"`
	BuGuMax          float64 `json:"buGuMax"`
	IbuMax           int     `json:"ibuMax"`
	ColorMin         int     `json:"colorMin"`
	OgMin            float64 `json:"ogMin"`
	Name             string  `json:"name"`
	CategoryNumber   string  `json:"categoryNumber"`
	BuGuMin          float64 `json:"buGuMin"`
	StyleGuide       string  `json:"styleGuide"`
	Type             string  `json:"type"`
	OgMax            float64 `json:"ogMax"`
	ColorMax         int     `json:"colorMax"`
	ID               string  `json:"_id"`
	StyleLetter      string  `json:"styleLetter"`
	FgMin            float64 `json:"fgMin"`
	IbuMin           int     `json:"ibuMin"`
	AbvMin           float64 `json:"abvMin"`
	FgMax            float64 `json:"fgMax"`
	RbrMin           float64 `json:"rbrMin"`
	LovibondMax      int     `json:"lovibondMax"`
	RbrMax           float64 `json:"rbrMax"`
	Category         string  `json:"category"`
	AbvMax           float64 `json:"abvMax"`
}

type CarbonationStyle struct {
	ID      string  `json:"_id"`
	Name    string  `json:"name"`
	CarbMax float64 `json:"carbMax"`
	CarbMin float64 `json:"carbMin"`
}

type MashFermentables struct {
	PotentialPercentage float64     `json:"potentialPercentage"`
	Notes               string      `json:"notes"`
	Color               int         `json:"color"`
	Potential           float64     `json:"potential"`
	BestBeforeDate      interface{} `json:"bestBeforeDate"`
	GrainCategory       string      `json:"grainCategory"`
	Percentage          float64     `json:"percentage"`
	CostPerAmount       interface{} `json:"costPerAmount"`
	ManufacturingDate   interface{} `json:"manufacturingDate"`
	ID                  string      `json:"_id"`
	Time                int         `json:"time,omitempty"`
	Supplier            string      `json:"supplier"`
	UsedIn              string      `json:"usedIn"`
	UserNotes           string      `json:"userNotes"`
	Amount              float64     `json:"amount"`
	DiastaticPower      interface{} `json:"diastaticPower"`
	Type                string      `json:"type"`
	NotFermentable      bool        `json:"notFermentable"`
	Protein             interface{} `json:"protein"`
	Attenuation         float64     `json:"attenuation"`
	Name                string      `json:"name"`
	Moisture            interface{} `json:"moisture"`
	Origin              string      `json:"origin"`
	Lovibond            float64     `json:"lovibond"`
	Use                 string      `json:"use,omitempty"`
	IbuPerAmount        interface{} `json:"ibuPerAmount"`
	Inventory           interface{} `json:"inventory"`
	Substitutes         string      `json:"substitutes"`
}

type Data struct {
	SpargeWaterAmount       float64            `json:"spargeWaterAmount"`
	MashFermentablesAmount  float64            `json:"mashFermentablesAmount"`
	StrikeTemp              interface{}        `json:"strikeTemp"`
	OtherFermentables       []interface{}      `json:"otherFermentables"`
	MashFermentables        []MashFermentables `json:"mashFermentables"`
	MashVolumeSurplus       int                `json:"mashVolumeSurplus"`
	TotalWaterAmount        float64            `json:"totalWaterAmount"`
	HltWaterAmount          float64            `json:"hltWaterAmount"`
	HopsAmount              float64            `json:"hopsAmount"`
	TopUpWater              int                `json:"topUpWater"`
	MashWaterAmount         float64            `json:"mashWaterAmount"`
	MashVolume              float64            `json:"mashVolume"`
	OtherFermentablesAmount int                `json:"otherFermentablesAmount"`
}

type Miscs struct {
	Unit            string      `json:"unit"`
	Use             string      `json:"use"`
	ID              string      `json:"_id"`
	WaterAdjustment bool        `json:"waterAdjustment"`
	Type            string      `json:"type"`
	TimeIsDays      bool        `json:"timeIsDays"`
	Name            string      `json:"name"`
	Time            interface{} `json:"time"`
	Concentration   interface{} `json:"concentration,omitempty"`
	Amount          float64     `json:"amount"`
}

type Equipment struct {
	MashWaterVolumeLimitEnabled bool        `json:"mashWaterVolumeLimitEnabled"`
	Rev                         string      `json:"_rev"`
	FermenterVolume             float64     `json:"fermenterVolume"`
	WaterCalculation            string      `json:"waterCalculation"`
	AromaHopUtilization         float64     `json:"aromaHopUtilization"`
	BoilSize                    float64     `json:"boilSize"`
	EfficiencyType              string      `json:"efficiencyType"`
	FermenterLoss               float64     `json:"fermenterLoss"`
	MashTunLoss                 float64     `json:"mashTunLoss"`
	CalcBoilVolume              bool        `json:"calcBoilVolume"`
	AltitudeAdjustment          bool        `json:"altitudeAdjustment"`
	BoilOffPerHr                float64     `json:"boilOffPerHr"`
	BottlingVolume              float64     `json:"bottlingVolume"`
	Version                     string      `json:"_version"`
	FermenterTopUp              interface{} `json:"fermenterTopUp"`
	Timestamp                   Timestamp   `json:"_timestamp"`
	SpargeWaterFormula          string      `json:"spargeWaterFormula"`
	FermenterLossEstimate       int         `json:"fermenterLossEstimate"`
	TimestampMs                 int64       `json:"_timestamp_ms"`
	BrewhouseEfficiency         float64     `json:"brewhouseEfficiency"`
	BoilTime                    int         `json:"boilTime"`
	CalcMashEfficiency          bool        `json:"calcMashEfficiency"`
	ID                          string      `json:"_id"`
	PostBoilKettleVol           float64     `json:"postBoilKettleVol"`
	HopUtilization              int         `json:"hopUtilization"`
	HopstandTemperature         float64     `json:"hopstandTemperature"`
	Hidden                      bool        `json:"hidden"`
	CalcAromaHopUtilization     bool        `json:"calcAromaHopUtilization"`
	Altitude                    float64     `json:"altitude"`
	Efficiency                  float64     `json:"efficiency"`
	EvaporationRate             float64     `json:"evaporationRate"`
	SpargeWaterOverflow         string      `json:"spargeWaterOverflow"`
	Created                     Created     `json:"_created"`
	TrubChillerLoss             float64     `json:"trubChillerLoss"`
	MashEfficiency              float64     `json:"mashEfficiency"`
	MashTunDeadSpace            float64     `json:"mashTunDeadSpace"`
	Name                        string      `json:"name"`
	FermenterVolumeBeforeTopUp  float64     `json:"fermenterVolumeBeforeTopUp"`
	BoilTemp                    float64     `json:"boilTemp"`
	MashWaterFormula            string      `json:"mashWaterFormula"`
	BatchSize                   float64     `json:"batchSize"`
}

type Calories struct {
	KJ      float64 `json:"kJ"`
	Total   float64 `json:"total"`
	Alcohol float64 `json:"alcohol"`
	Carbs   float64 `json:"carbs"`
}

type Carbs struct {
	Total float64 `json:"total"`
}

type Nutrition struct {
	Calories Calories `json:"calories"`
	Carbs    Carbs    `json:"carbs"`
}
