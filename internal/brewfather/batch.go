//nolint:tagliatelle
package brewfather

type Batch struct {
	BoilSteps                     []BoilSteps              `json:"boilSteps"`
	MeasuredConversionEfficiency  interface{}              `json:"measuredConversionEfficiency"`
	Brewer                        interface{}              `json:"brewer"`
	BrewDate                      int64                    `json:"brewDate"`
	Recipe                        Recipe                   `json:"recipe"`
	BatchYeasts                   []interface{}            `json:"batchYeasts"`
	Devices                       Devices                  `json:"devices"`
	EstimatedOg                   float64                  `json:"estimatedOg"`
	HideBatchEvents               bool                     `json:"hideBatchEvents"`
	Events                        []Events                 `json:"events"`
	EstimatedRbRatio              float64                  `json:"estimatedRbRatio"`
	MeasuredMashEfficiency        float64                  `json:"measuredMashEfficiency"`
	BatchMiscs                    []interface{}            `json:"batchMiscs"`
	EstimatedBuGuRatio            float64                  `json:"estimatedBuGuRatio"`
	Created                       Created                  `json:"_created"`
	MeasuredKettleEfficiency      float64                  `json:"measuredKettleEfficiency"`
	Status                        string                   `json:"status"`
	Version                       string                   `json:"_version"`
	Notes                         []Notes                  `json:"notes"`
	EstimatedTotalGravity         float64                  `json:"estimatedTotalGravity"`
	PrimingSugarEquiv             interface{}              `json:"primingSugarEquiv"`
	Cost                          Cost                     `json:"cost"`
	FermentationControllerEnabled bool                     `json:"fermentationControllerEnabled"`
	MeasuredAbv                   float64                  `json:"measuredAbv"`
	BatchMiscsLocal               []interface{}            `json:"batchMiscsLocal"`
	Rev                           string                   `json:"_rev"`
	Init                          bool                     `json:"_init"`
	Name                          string                   `json:"name"`
	CarbonationForce              float64                  `json:"carbonationForce"`
	CarbonationType               string                   `json:"carbonationType"`
	MashStepsCount                int                      `json:"mashStepsCount"`
	BatchFermentables             []BatchFermentables      `json:"batchFermentables"`
	BoilStepsCount                int                      `json:"boilStepsCount"`
	Measurements                  []interface{}            `json:"measurements"`
	BatchHopsLocal                []interface{}            `json:"batchHopsLocal"`
	MeasuredAttenuation           float64                  `json:"measuredAttenuation"`
	Timestamp                     Timestamp                `json:"_timestamp"`
	TimestampMs                   int64                    `json:"_timestamp_ms"`
	Type                          string                   `json:"_type"`
	Archived                      bool                     `json:"_archived"`
	BrewControllerEnabled         bool                     `json:"brewControllerEnabled"`
	BatchFermentablesLocal        []BatchFermentablesLocal `json:"batchFermentablesLocal"`
	EstimatedIbu                  int                      `json:"estimatedIbu"`
	BatchNo                       int                      `json:"batchNo"`
	FermentationStartDate         int64                    `json:"fermentationStartDate"`
	EstimatedColor                float64                  `json:"estimatedColor"`
	BottlingDate                  int64                    `json:"bottlingDate"`
	BatchYeastsLocal              []BatchYeastsLocal       `json:"batchYeastsLocal"`
	BatchHops                     []BatchHops              `json:"batchHops"`
	EstimatedFg                   float64                  `json:"estimatedFg"`
	Hidden                        bool                     `json:"hidden"`
	ID                            string                   `json:"_id"`
	MeasuredEfficiency            float64                  `json:"measuredEfficiency"`
}

type BoilSteps struct {
	Time int    `json:"time"`
	Name string `json:"name"`
}

type Style struct {
	Category       string      `json:"category"`
	Mouthfeel      interface{} `json:"mouthfeel"`
	RbrMax         float64     `json:"rbrMax"`
	Type           string      `json:"type"`
	Profile        string      `json:"profile"`
	Origin         interface{} `json:"origin"`
	IbuMax         int         `json:"ibuMax"`
	CategoryNumber int         `json:"categoryNumber"`
	IbuMin         int         `json:"ibuMin"`
	RbrMin         float64     `json:"rbrMin"`
	StyleGuide     string      `json:"styleGuide"`
	Ingredients    string      `json:"ingredients"`
	ColorMax       int         `json:"colorMax"`
	AbvMin         int         `json:"abvMin"`
	FgMax          float64     `json:"fgMax"`
	LovibondMin    int         `json:"lovibondMin"`
	AbvMax         float64     `json:"abvMax"`
	Flavor         interface{} `json:"flavor"`
	BuGuMax        float64     `json:"buGuMax"`
	OgMax          float64     `json:"ogMax"`
	FgMin          float64     `json:"fgMin"`
	Appearance     interface{} `json:"appearance"`
	History        interface{} `json:"history"`
	ColorMin       int         `json:"colorMin"`
	OgMin          float64     `json:"ogMin"`
	Name           string      `json:"name"`
	LovibondMax    int         `json:"lovibondMax"`
	Notes          string      `json:"notes"`
	BuGuMin        float64     `json:"buGuMin"`
	StyleLetter    string      `json:"styleLetter"`
	Aroma          interface{} `json:"aroma"`
	Examples       string      `json:"examples"`
}

type Timestamp struct {
	Seconds     int `json:"_seconds"`
	Nanoseconds int `json:"_nanoseconds"`
}

type Created struct {
	Seconds     int `json:"_seconds"`
	Nanoseconds int `json:"_nanoseconds"`
}

type Source struct {
	Alkalinity                 float64   `json:"alkalinity"`
	ResidualAlkalinityMeqLCalc float64   `json:"residualAlkalinityMeqLCalc"`
	ResidualAlkalinity         float64   `json:"residualAlkalinity"`
	Hardness                   int       `json:"hardness"`
	Magnesium                  float64   `json:"magnesium"`
	Version                    string    `json:"_version"`
	Name                       string    `json:"name"`
	SoClRatio                  float64   `json:"soClRatio"`
	Bicarbonate                int       `json:"bicarbonate"`
	ID                         string    `json:"_id"`
	Rev                        string    `json:"_rev"`
	Anions                     float64   `json:"anions"`
	IonBalanceOff              bool      `json:"ionBalanceOff"`
	Sodium                     int       `json:"sodium"`
	TimestampMs                int64     `json:"_timestamp_ms"`
	Chloride                   float64   `json:"chloride"`
	Sulfate                    float64   `json:"sulfate"`
	Type                       string    `json:"type"`
	Ph                         float64   `json:"ph"`
	Calcium                    float64   `json:"calcium"`
	Cations                    float64   `json:"cations"`
	Hidden                     bool      `json:"hidden"`
	Timestamp                  Timestamp `json:"_timestamp"`
	Created                    Created   `json:"_created"`
	BicarbonateMeqL            float64   `json:"bicarbonateMeqL"`
	IonBalance                 int       `json:"ionBalance"`
}

type Total struct {
	Anions                     float64   `json:"anions"`
	ResidualAlkalinityMeqLCalc float64   `json:"residualAlkalinityMeqLCalc"`
	Chloride                   float64   `json:"chloride"`
	Cations                    float64   `json:"cations"`
	Ph                         float64   `json:"ph"`
	TimestampMs                int64     `json:"_timestamp_ms"`
	Alkalinity                 float64   `json:"alkalinity"`
	Sulfate                    float64   `json:"sulfate"`
	BicarbonateMeqL            float64   `json:"bicarbonateMeqL"`
	SoClRatio                  float64   `json:"soClRatio"`
	Created                    Created   `json:"_created"`
	Type                       string    `json:"type"`
	Sodium                     int       `json:"sodium"`
	Bicarbonate                int       `json:"bicarbonate"`
	Calcium                    float64   `json:"calcium"`
	Magnesium                  float64   `json:"magnesium"`
	Name                       string    `json:"name"`
	IonBalanceOff              bool      `json:"ionBalanceOff"`
	Rev                        string    `json:"_rev"`
	Version                    string    `json:"_version"`
	ResidualAlkalinity         float64   `json:"residualAlkalinity"`
	Hidden                     bool      `json:"hidden"`
	ID                         string    `json:"_id"`
	IonBalance                 int       `json:"ionBalance"`
	Timestamp                  Timestamp `json:"_timestamp"`
	Hardness                   int       `json:"hardness"`
}

type MashAdjustmentsAcids struct {
	AlkalinityMeqL int    `json:"alkalinityMeqL"`
	Type           string `json:"type"`
	Concentration  int    `json:"concentration"`
	Amount         int    `json:"amount"`
}

type MashAdjustments struct {
	CalciumHydroxide       interface{}            `json:"calciumHydroxide"`
	SodiumChloride         int                    `json:"sodiumChloride"`
	Acids                  []MashAdjustmentsAcids `json:"acids"`
	SodiumMetabisulfite    int                    `json:"sodiumMetabisulfite"`
	LtAMS                  int                    `json:"ltAMS"`
	CalciumCarbonate       int                    `json:"calciumCarbonate"`
	Magnesium              float64                `json:"magnesium"`
	LtDWB                  int                    `json:"ltDWB"`
	Chloride               float64                `json:"chloride"`
	MagnesiumSulfate       interface{}            `json:"magnesiumSulfate"`
	Sodium                 int                    `json:"sodium"`
	MagnesiumChloride      int                    `json:"magnesiumChloride"`
	SodiumBicarbonate      interface{}            `json:"sodiumBicarbonate"`
	Calcium                float64                `json:"calcium"`
	CalciumSulfate         interface{}            `json:"calciumSulfate"`
	CalciumChloride        interface{}            `json:"calciumChloride"`
	SodiumMetabisulfitePPM int                    `json:"sodiumMetabisulfitePPM"`
	Bicarbonate            int                    `json:"bicarbonate"`
	Volume                 float64                `json:"volume"`
	Sulfate                float64                `json:"sulfate"`
}

type SpargeAdjustmentsAcids struct {
	Concentration int    `json:"concentration"`
	Amount        int    `json:"amount"`
	Type          string `json:"type"`
}

type SpargeAdjustments struct {
	Sodium                 int                      `json:"sodium"`
	MagnesiumChloride      int                      `json:"magnesiumChloride"`
	Calcium                float64                  `json:"calcium"`
	SodiumBicarbonate      int                      `json:"sodiumBicarbonate"`
	MagnesiumSulfate       float64                  `json:"magnesiumSulfate"`
	CalciumSulfate         float64                  `json:"calciumSulfate"`
	LtAMS                  int                      `json:"ltAMS"`
	CalciumCarbonate       int                      `json:"calciumCarbonate"`
	SodiumMetabisulfite    int                      `json:"sodiumMetabisulfite"`
	SodiumChloride         int                      `json:"sodiumChloride"`
	CalciumHydroxide       int                      `json:"calciumHydroxide"`
	Magnesium              float64                  `json:"magnesium"`
	Acids                  []SpargeAdjustmentsAcids `json:"acids"`
	Volume                 float64                  `json:"volume"`
	Bicarbonate            int                      `json:"bicarbonate"`
	Chloride               float64                  `json:"chloride"`
	SodiumMetabisulfitePPM int                      `json:"sodiumMetabisulfitePPM"`
	CalciumChloride        float64                  `json:"calciumChloride"`
	Sulfate                float64                  `json:"sulfate"`
	LtDWB                  int                      `json:"ltDWB"`
}

type Sparge struct {
	Hidden                     bool      `json:"hidden"`
	ResidualAlkalinityMeqLCalc float64   `json:"residualAlkalinityMeqLCalc"`
	Rev                        string    `json:"_rev"`
	Created                    Created   `json:"_created"`
	Sulfate                    float64   `json:"sulfate"`
	TimestampMs                int64     `json:"_timestamp_ms"`
	Type                       string    `json:"type"`
	IonBalanceOff              bool      `json:"ionBalanceOff"`
	BicarbonateMeqL            float64   `json:"bicarbonateMeqL"`
	SoClRatio                  float64   `json:"soClRatio"`
	Name                       string    `json:"name"`
	Alkalinity                 float64   `json:"alkalinity"`
	Magnesium                  float64   `json:"magnesium"`
	Chloride                   float64   `json:"chloride"`
	ID                         string    `json:"_id"`
	Bicarbonate                int       `json:"bicarbonate"`
	Calcium                    float64   `json:"calcium"`
	Version                    string    `json:"_version"`
	Hardness                   int       `json:"hardness"`
	Anions                     float64   `json:"anions"`
	ResidualAlkalinity         float64   `json:"residualAlkalinity"`
	Cations                    float64   `json:"cations"`
	Ph                         float64   `json:"ph"`
	Timestamp                  Timestamp `json:"_timestamp"`
	Sodium                     int       `json:"sodium"`
	IonBalance                 int       `json:"ionBalance"`
}

type Meta struct {
	EqualSourceTotal bool `json:"equalSourceTotal"`
}

type WaterMash struct {
	Sulfate                    float64   `json:"sulfate"`
	Anions                     float64   `json:"anions"`
	ID                         string    `json:"_id"`
	Name                       string    `json:"name"`
	IonBalance                 int       `json:"ionBalance"`
	Sodium                     int       `json:"sodium"`
	BicarbonateMeqL            float64   `json:"bicarbonateMeqL"`
	SoClRatio                  float64   `json:"soClRatio"`
	Bicarbonate                int       `json:"bicarbonate"`
	IonBalanceOff              bool      `json:"ionBalanceOff"`
	Rev                        string    `json:"_rev"`
	ResidualAlkalinity         float64   `json:"residualAlkalinity"`
	Calcium                    float64   `json:"calcium"`
	TimestampMs                int64     `json:"_timestamp_ms"`
	Magnesium                  float64   `json:"magnesium"`
	ResidualAlkalinityMeqLCalc float64   `json:"residualAlkalinityMeqLCalc"`
	Ph                         float64   `json:"ph"`
	Type                       string    `json:"type"`
	Timestamp                  Timestamp `json:"_timestamp"`
	Alkalinity                 float64   `json:"alkalinity"`
	Hidden                     bool      `json:"hidden"`
	Cations                    float64   `json:"cations"`
	Chloride                   float64   `json:"chloride"`
	Created                    Created   `json:"_created"`
	Hardness                   int       `json:"hardness"`
	Version                    string    `json:"_version"`
}

type SodiumBicarbonate struct {
	Mash   bool `json:"mash"`
	Auto   bool `json:"auto"`
	Sparge bool `json:"sparge"`
}

type MagnesiumSulfate struct {
	Auto   bool `json:"auto"`
	Mash   bool `json:"mash"`
	Sparge bool `json:"sparge"`
}

type CalciumChloride struct {
	Auto   bool   `json:"auto"`
	Mash   bool   `json:"mash"`
	Form   string `json:"form"`
	Sparge bool   `json:"sparge"`
}

type CalciumHydroxide struct {
	Auto   bool `json:"auto"`
	Mash   bool `json:"mash"`
	Sparge bool `json:"sparge"`
}

type CalciumSulfate struct {
	Mash   bool `json:"mash"`
	Sparge bool `json:"sparge"`
	Auto   bool `json:"auto"`
}

type Settings struct {
	SodiumBicarbonate SodiumBicarbonate `json:"sodiumBicarbonate"`
	MagnesiumSulfate  MagnesiumSulfate  `json:"magnesiumSulfate"`
	AdjustSparge      bool              `json:"adjustSparge"`
	CalciumChloride   CalciumChloride   `json:"calciumChloride"`
	CalciumHydroxide  CalciumHydroxide  `json:"calciumHydroxide"`
	CalciumSulfate    CalciumSulfate    `json:"calciumSulfate"`
}

type TotalAdjustments struct {
	Chloride               float64 `json:"chloride"`
	Calcium                float64 `json:"calcium"`
	SodiumMetabisulfite    int     `json:"sodiumMetabisulfite"`
	LtAMS                  int     `json:"ltAMS"`
	MagnesiumSulfate       float64 `json:"magnesiumSulfate"`
	Sodium                 int     `json:"sodium"`
	Bicarbonate            int     `json:"bicarbonate"`
	CalciumHydroxide       int     `json:"calciumHydroxide"`
	MagnesiumChloride      int     `json:"magnesiumChloride"`
	CalciumChloride        float64 `json:"calciumChloride"`
	SodiumBicarbonate      int     `json:"sodiumBicarbonate"`
	Sulfate                float64 `json:"sulfate"`
	SodiumMetabisulfitePPM int     `json:"sodiumMetabisulfitePPM"`
	CalciumSulfate         float64 `json:"calciumSulfate"`
	SodiumChloride         int     `json:"sodiumChloride"`
	Volume                 float64 `json:"volume"`
	Magnesium              float64 `json:"magnesium"`
	CalciumCarbonate       int     `json:"calciumCarbonate"`
	LtDWB                  int     `json:"ltDWB"`
}

type Water struct {
	SpargeWaterAmount       interface{}       `json:"spargeWaterAmount"`
	Source                  Source            `json:"source"`
	Total                   Total             `json:"total"`
	MashPh                  float64           `json:"mashPh"`
	DilutionPercentage      interface{}       `json:"dilutionPercentage"`
	MashWaterAmount         interface{}       `json:"mashWaterAmount"`
	SpargeAcidPhAdjustment  int               `json:"spargeAcidPhAdjustment"`
	MashAdjustments         MashAdjustments   `json:"mashAdjustments"`
	MashTargetDiff          interface{}       `json:"mashTargetDiff"`
	SpargeAdjustments       SpargeAdjustments `json:"spargeAdjustments"`
	AcidPhAdjustment        int               `json:"acidPhAdjustment"`
	SpargeTargetDiff        interface{}       `json:"spargeTargetDiff"`
	SourceTargetDiff        interface{}       `json:"sourceTargetDiff"`
	Sparge                  Sparge            `json:"sparge"`
	Meta                    Meta              `json:"meta"`
	Mash                    WaterMash         `json:"mash"`
	Settings                Settings          `json:"settings"`
	Diluted                 interface{}       `json:"diluted"`
	TotalTargetDiff         interface{}       `json:"totalTargetDiff"`
	TotalAdjustments        TotalAdjustments  `json:"totalAdjustments"`
	MashPhDistilled         float64           `json:"mashPhDistilled"`
	EnableSpargeAdjustments bool              `json:"enableSpargeAdjustments"`
}

type Equipment struct {
	SpargeTemperature           string      `json:"spargeTemperature"`
	PostBoilKettleVol           float64     `json:"postBoilKettleVol"`
	Efficiency                  float64     `json:"efficiency"`
	EfficiencyType              string      `json:"efficiencyType"`
	EvaporationRate             float64     `json:"evaporationRate"`
	MashTunDeadSpace            float64     `json:"mashTunDeadSpace"`
	CalcMashEfficiency          bool        `json:"calcMashEfficiency"`
	CalcBoilVolume              bool        `json:"calcBoilVolume"`
	TrubChillerLoss             float64     `json:"trubChillerLoss"`
	FermenterLoss               interface{} `json:"fermenterLoss"`
	ID                          interface{} `json:"_id"`
	BottlingVolume              float64     `json:"bottlingVolume"`
	AromaHopUtilization         float64     `json:"aromaHopUtilization"`
	MashWaterVolumeLimitEnabled bool        `json:"mashWaterVolumeLimitEnabled"`
	MashEfficiency              float64     `json:"mashEfficiency"`
	HopUtilization              int         `json:"hopUtilization"`
	MashWaterFormula            interface{} `json:"mashWaterFormula"`
	BoilOffPerHr                float64     `json:"boilOffPerHr"`
	MashWaterMax                float64     `json:"mashWaterMax"`
	BatchSize                   float64     `json:"batchSize"`
	FermenterLossEstimate       int         `json:"fermenterLossEstimate"`
	HopstandTemperature         float64     `json:"hopstandTemperature"`
	BoilSize                    float64     `json:"boilSize"`
	SpargeWaterFormula          interface{} `json:"spargeWaterFormula"`
	FermenterVolume             float64     `json:"fermenterVolume"`
	CalcAromaHopUtilization     bool        `json:"calcAromaHopUtilization"`
	BoilTime                    int         `json:"boilTime"`
	Name                        string      `json:"name"`
	AmbientTemperature          string      `json:"ambientTemperature"`
	GrainTemperature            string      `json:"grainTemperature"`
}

type MashFermentables struct {
	IbuPerAmount        int         `json:"ibuPerAmount"`
	AddAfterBoil        string      `json:"addAfterBoil"`
	Name                string      `json:"name"`
	Type                string      `json:"type"`
	Origin              string      `json:"origin"`
	GrainCategory       string      `json:"grainCategory"`
	Percentage          float64     `json:"percentage"`
	Supplier            string      `json:"supplier"`
	ID                  string      `json:"_id"`
	NotFermentable      interface{} `json:"notFermentable"`
	Potential           float64     `json:"potential"`
	PotentialPercentage float64     `json:"potentialPercentage"`
	Color               float64     `json:"color"`
	Notes               string      `json:"notes"`
	Amount              float64     `json:"amount"`
}

type OtherFermentables struct {
	ID                  string      `json:"_id"`
	PotentialPercentage float64     `json:"potentialPercentage"`
	Supplier            string      `json:"supplier"`
	Potential           float64     `json:"potential"`
	Notes               string      `json:"notes"`
	Amount              float64     `json:"amount"`
	NotFermentable      interface{} `json:"notFermentable"`
	Origin              string      `json:"origin"`
	GrainCategory       interface{} `json:"grainCategory"`
	Type                string      `json:"type"`
	Color               float64     `json:"color"`
	AddAfterBoil        string      `json:"addAfterBoil"`
	IbuPerAmount        int         `json:"ibuPerAmount"`
	Percentage          float64     `json:"percentage"`
	Name                string      `json:"name"`
}

type Data struct {
	BatchSpargeWaterAmount1 interface{}         `json:"batchSpargeWaterAmount1"`
	BatchSpargeWaterAmount3 interface{}         `json:"batchSpargeWaterAmount3"`
	StrikeTemp              interface{}         `json:"strikeTemp"`
	AllDiastaticPower       bool                `json:"allDiastaticPower"`
	BatchSpargeWaterAmount4 interface{}         `json:"batchSpargeWaterAmount4"`
	MashWaterAmount         float64             `json:"mashWaterAmount"`
	BatchSpargeWaterAmount2 interface{}         `json:"batchSpargeWaterAmount2"`
	TopUpWater              int                 `json:"topUpWater"`
	OtherFermentablesAmount float64             `json:"otherFermentablesAmount"`
	MashFermentables        []MashFermentables  `json:"mashFermentables"`
	MashVolumeSurplus       int                 `json:"mashVolumeSurplus"`
	HltWaterAmount          float64             `json:"hltWaterAmount"`
	SpargeWaterAmount       float64             `json:"spargeWaterAmount"`
	TotalWaterAmount        float64             `json:"totalWaterAmount"`
	OtherFermentables       []OtherFermentables `json:"otherFermentables"`
	MashFermentablesAmount  float64             `json:"mashFermentablesAmount"`
	TotalDiastaticPower     int                 `json:"totalDiastaticPower"`
	MashVolume              float64             `json:"mashVolume"`
	HopsAmount              float64             `json:"hopsAmount"`
}

type Yeasts struct {
	Unit        string  `json:"unit"`
	Form        string  `json:"form"`
	Notes       string  `json:"notes"`
	Amount      float64 `json:"amount"`
	Attenuation int     `json:"attenuation"`
	ProductID   string  `json:"productId"`
	Laboratory  string  `json:"laboratory"`
	ID          string  `json:"_id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
}

type Steps struct {
	RampTime int     `json:"rampTime"`
	Name     string  `json:"name"`
	StepTime int     `json:"stepTime"`
	StepTemp float64 `json:"stepTemp"`
	Type     string  `json:"type"`
}

type RecipeMash struct {
	ID    interface{} `json:"_id"`
	Name  string      `json:"name"`
	Steps []Steps     `json:"steps"`
}

type Defaults struct {
	Preferred   string `json:"preferred"`
	Ibu         string `json:"ibu"`
	Color       string `json:"color"`
	Hop         string `json:"hop"`
	Attenuation string `json:"attenuation"`
	Altitude    string `json:"altitude"`
	Abv         string `json:"abv"`
	GrainColor  string `json:"grainColor"`
	Volume      string `json:"volume"`
	Temp        string `json:"temp"`
	Pressure    string `json:"pressure"`
	Gravity     string `json:"gravity"`
	Carbonation string `json:"carbonation"`
	Weight      string `json:"weight"`
}

type Fermentables struct {
	Potential           float64     `json:"potential"`
	PotentialPercentage float64     `json:"potentialPercentage"`
	GrainCategory       string      `json:"grainCategory"`
	Amount              float64     `json:"amount"`
	Notes               string      `json:"notes"`
	Supplier            string      `json:"supplier"`
	Percentage          float64     `json:"percentage"`
	Name                string      `json:"name"`
	IbuPerAmount        int         `json:"ibuPerAmount"`
	NotFermentable      interface{} `json:"notFermentable"`
	AddAfterBoil        string      `json:"addAfterBoil"`
	Color               float64     `json:"color"`
	ID                  string      `json:"_id"`
	Origin              string      `json:"origin"`
	Type                string      `json:"type"`
}

type Carbs struct {
	Total float64 `json:"total"`
}

type Calories struct {
	KJ      float64 `json:"kJ"`
	Total   float64 `json:"total"`
	Alcohol float64 `json:"alcohol"`
	Carbs   float64 `json:"carbs"`
}

type Nutrition struct {
	Carbs    Carbs    `json:"carbs"`
	Calories Calories `json:"calories"`
}

type FermentationStep struct {
	Type       string  `json:"type"`
	ActualTime int64   `json:"actualTime"`
	StepTemp   float64 `json:"stepTemp"`
	StepTime   int     `json:"stepTime"`
}

type Fermentation struct {
	ID    interface{}        `json:"_id"`
	Name  string             `json:"name"`
	Steps []FermentationStep `json:"steps"`
}

type CarbonationStyle struct {
	ID      string  `json:"_id"`
	Name    string  `json:"name"`
	CarbMax float64 `json:"carbMax"`
	CarbMin float64 `json:"carbMin"`
}

type Hops struct {
	Temp   interface{} `json:"temp"`
	Notes  string      `json:"notes"`
	Name   string      `json:"name"`
	Alpha  float64     `json:"alpha"`
	Use    string      `json:"use"`
	Ibu    float64     `json:"ibu"`
	Type   string      `json:"type"`
	Amount float64     `json:"amount"`
	Origin string      `json:"origin"`
	ID     string      `json:"_id"`
	Time   int         `json:"time"`
}

type Recipe struct {
	Style                    Style            `json:"style"`
	Rev                      string           `json:"_rev"`
	Hidden                   bool             `json:"hidden"`
	ManualFg                 bool             `json:"manualFg"`
	SearchTags               []string         `json:"searchTags"`
	FermentablesTotalAmount  float64          `json:"fermentablesTotalAmount"`
	HopStandMinutes          int              `json:"hopStandMinutes"`
	BoilTime                 int              `json:"boilTime"`
	Water                    Water            `json:"water"`
	Init                     bool             `json:"_init"`
	Equipment                Equipment        `json:"equipment"`
	Created                  Created          `json:"_created"`
	Ev                       float64          `json:"_ev"`
	StyleOg                  bool             `json:"styleOg"`
	Attenuation              float64          `json:"attenuation"`
	MashEfficiency           float64          `json:"mashEfficiency"`
	Fg                       float64          `json:"fg"`
	PostBoilGravity          float64          `json:"postBoilGravity"`
	Data                     Data             `json:"data"`
	BatchSize                float64          `json:"batchSize"`
	RbRatio                  float64          `json:"rbRatio"`
	Carbonation              float64          `json:"carbonation"`
	StyleAbv                 bool             `json:"styleAbv"`
	Author                   string           `json:"author"`
	StyleConformity          bool             `json:"styleConformity"`
	Efficiency               float64          `json:"efficiency"`
	Yeasts                   []Yeasts         `json:"yeasts"`
	Name                     string           `json:"name"`
	PrimaryTemp              float64          `json:"primaryTemp"`
	Mash                     RecipeMash       `json:"mash"`
	Notes                    string           `json:"notes"`
	StyleRbr                 bool             `json:"styleRbr"`
	FgFormula                string           `json:"fgFormula"`
	ExtraGravity             int              `json:"extraGravity"`
	BoilSize                 float64          `json:"boilSize"`
	StyleColor               bool             `json:"styleColor"`
	Defaults                 Defaults         `json:"defaults"`
	TimestampMs              int64            `json:"_timestamp_ms"`
	Path                     string           `json:"path"`
	Public                   bool             `json:"public"`
	Fermentables             []Fermentables   `json:"fermentables"`
	Og                       float64          `json:"og"`
	Nutrition                Nutrition        `json:"nutrition"`
	BuGuRatio                float64          `json:"buGuRatio"`
	OgPlato                  float64          `json:"ogPlato"`
	StyleCarb                bool             `json:"styleCarb"`
	StyleIbu                 bool             `json:"styleIbu"`
	DiastaticPower           int              `json:"diastaticPower"`
	FirstWortGravity         interface{}      `json:"firstWortGravity"`
	Fermentation             Fermentation     `json:"fermentation"`
	FgEstimated              float64          `json:"fgEstimated"`
	Ibu                      float64          `json:"ibu"`
	Type                     string           `json:"type"`
	Color                    float64          `json:"color"`
	Timestamp                Timestamp        `json:"_timestamp"`
	CarbonationStyle         CarbonationStyle `json:"carbonationStyle"`
	Miscs                    []interface{}    `json:"miscs"`
	AvgWeightedHopstandTemp  float64          `json:"avgWeightedHopstandTemp"`
	StyleFg                  bool             `json:"styleFg"`
	IbuFormula               string           `json:"ibuFormula"`
	Share                    interface{}      `json:"_share"`
	StyleBuGu                bool             `json:"styleBuGu"`
	ID                       string           `json:"_id"`
	FermentableIbu           int              `json:"fermentableIbu"`
	TotalGravity             float64          `json:"totalGravity"`
	HopsTotalAmount          float64          `json:"hopsTotalAmount"`
	Tags                     interface{}      `json:"tags"`
	Hops                     []Hops           `json:"hops"`
	Version                  string           `json:"_version"`
	PreBoilGravity           float64          `json:"preBoilGravity"`
	YeastToleranceExceededBy interface{}      `json:"yeastToleranceExceededBy"`
	Abv                      float64          `json:"abv"`
	SumDryHopPerLiter        int              `json:"sumDryHopPerLiter"`
}

type FloatyHydrometer struct {
	Items   []interface{} `json:"items"`
	Enabled bool          `json:"enabled"`
}

type MyBrewbot struct {
	Enabled bool          `json:"enabled"`
	Items   []interface{} `json:"items"`
}

type FloatHydrometer struct {
	Enabled bool          `json:"enabled"`
	Items   []interface{} `json:"items"`
}

type Gfcc struct {
	BrewDeviceID interface{}   `json:"brewDeviceId"`
	Enabled      bool          `json:"enabled"`
	Items        []interface{} `json:"items"`
}

type BrewPiLess struct {
	Enabled bool          `json:"enabled"`
	Items   []interface{} `json:"items"`
}

type SmartPid struct {
	BrewDeviceID interface{}   `json:"brewDeviceId"`
	Items        []interface{} `json:"items"`
	Enabled      bool          `json:"enabled"`
}

type ISpindel struct {
	Enabled bool          `json:"enabled"`
	Items   []interface{} `json:"items"`
}

type Stream struct {
	Items   []interface{} `json:"items"`
	Enabled bool          `json:"enabled"`
}

type PlaatoKeg struct {
	Items   []interface{} `json:"items"`
	Enabled bool          `json:"enabled"`
}

type PlaatoAirlock struct {
	Enabled bool          `json:"enabled"`
	Items   []interface{} `json:"items"`
}

type Tilt struct {
	Enabled bool          `json:"enabled"`
	Gravity bool          `json:"gravity"`
	Temp    bool          `json:"temp"`
	Mode    string        `json:"mode"`
	Items   []interface{} `json:"items"`
}

type Devices struct {
	FloatyHydrometer FloatyHydrometer `json:"floatyHydrometer"`
	MyBrewbot        MyBrewbot        `json:"myBrewbot"`
	FloatHydrometer  FloatHydrometer  `json:"floatHydrometer"`
	Gfcc             Gfcc             `json:"gfcc"`
	BrewPiLess       BrewPiLess       `json:"brewPiLess"`
	SmartPid         SmartPid         `json:"smartPid"`
	ISpindel         ISpindel         `json:"iSpindel"`
	Stream           Stream           `json:"stream"`
	PlaatoKeg        PlaatoKeg        `json:"plaatoKeg"`
	PlaatoAirlock    PlaatoAirlock    `json:"plaatoAirlock"`
	Tilt             Tilt             `json:"tilt"`
}

type Events struct {
	EventText       string `json:"eventText"`
	EventType       string `json:"eventType"`
	Title           string `json:"title"`
	Time            int64  `json:"time"`
	DayEvent        bool   `json:"dayEvent"`
	Description     string `json:"description"`
	Active          bool   `json:"active"`
	DescriptionHTML string `json:"descriptionHTML"`
	NotifyTime      int64  `json:"notifyTime,omitempty"`
}

type Notes struct {
	Timestamp int64  `json:"timestamp"`
	Status    string `json:"status"`
	Type      string `json:"type"`
	Note      string `json:"note"`
}

type Cost struct {
	Fermentables      int `json:"fermentables"`
	Hops              int `json:"hops"`
	Miscs             int `json:"miscs"`
	Total             int `json:"total"`
	YeastsShare       int `json:"yeastsShare"`
	Yeasts            int `json:"yeasts"`
	FermentablesShare int `json:"fermentablesShare"`
	MiscsShare        int `json:"miscsShare"`
	PerBottlingLiter  int `json:"perBottlingLiter"`
	HopsShare         int `json:"hopsShare"`
}

type BatchFermentables struct {
	TotalCost           int         `json:"totalCost"`
	Origin              string      `json:"origin"`
	Name                string      `json:"name"`
	NotFermentable      interface{} `json:"notFermentable"`
	Inventory           int         `json:"inventory"`
	PotentialPercentage float64     `json:"potentialPercentage"`
	NotInRecipe         bool        `json:"notInRecipe"`
	Type                string      `json:"type"`
	DisplayAmount       float64     `json:"displayAmount"`
	Amount              float64     `json:"amount"`
	Potential           float64     `json:"potential"`
	CostPerAmount       int         `json:"costPerAmount"`
	Attenuation         interface{} `json:"attenuation"`
	ID                  string      `json:"_id"`
	Color               float64     `json:"color"`
	IbuPerAmount        interface{} `json:"ibuPerAmount"`
	Supplier            string      `json:"supplier"`
}

type BatchFermentablesLocal struct {
	Color               float64     `json:"color"`
	AddAfterBoil        string      `json:"addAfterBoil"`
	CostPerAmount       int         `json:"costPerAmount"`
	NotFermentable      interface{} `json:"notFermentable"`
	Type                string      `json:"type"`
	Potential           float64     `json:"potential"`
	GrainCategory       string      `json:"grainCategory"`
	DisplayAmount       float64     `json:"displayAmount"`
	Percentage          float64     `json:"percentage"`
	Notes               string      `json:"notes"`
	ID                  string      `json:"_id"`
	Amount              float64     `json:"amount"`
	Origin              string      `json:"origin"`
	IbuPerAmount        int         `json:"ibuPerAmount"`
	NotInRecipe         bool        `json:"notInRecipe"`
	PotentialPercentage int         `json:"potentialPercentage"`
	Inventory           int         `json:"inventory"`
	Supplier            string      `json:"supplier"`
	Name                string      `json:"name"`
}

type BatchYeastsLocal struct {
	Amount        float64 `json:"amount"`
	Attenuation   int     `json:"attenuation"`
	Unit          string  `json:"unit"`
	Name          string  `json:"name"`
	Laboratory    string  `json:"laboratory"`
	CostPerAmount int     `json:"costPerAmount"`
	ProductID     string  `json:"productId"`
	Inventory     int     `json:"inventory"`
	Form          string  `json:"form"`
	NotInRecipe   bool    `json:"notInRecipe"`
	Type          string  `json:"type"`
	Notes         string  `json:"notes"`
	InventoryUnit string  `json:"inventoryUnit"`
	ID            string  `json:"_id"`
	DisplayAmount float64 `json:"displayAmount"`
}

type BatchHops struct {
	Type          string  `json:"type"`
	Alpha         float64 `json:"alpha"`
	CostPerAmount int     `json:"costPerAmount"`
	Origin        string  `json:"origin"`
	Usage         string  `json:"usage"`
	ID            string  `json:"_id"`
	Inventory     int     `json:"inventory"`
	TotalCost     int     `json:"totalCost"`
	Amount        float64 `json:"amount"`
	NotInRecipe   bool    `json:"notInRecipe"`
	Name          string  `json:"name"`
	DisplayAmount float64 `json:"displayAmount"`
}
