//nolint:tagliatelle // this file was manual generated and includes minor changes
package brewfather

type BatchSummary struct {
	ID       string        `json:"_id"`
	Name     string        `json:"name"`
	BatchNo  int           `json:"batchNo"`
	Status   string        `json:"status"`
	Brewer   interface{}   `json:"brewer"`
	BrewDate int64         `json:"brewDate"`
	Recipe   RecipeSummary `json:"recipe"`
}

type RecipeSummary struct {
	Name string `json:"name"`
}

type BatchDetail struct {
	Name                          string                   `json:"name"`
	MeasuredBatchSize             float64                  `json:"measuredBatchSize"`
	Archived                      bool                     `json:"_archived"`
	BoilSteps                     []BoilSteps              `json:"boilSteps"`
	BatchFermentablesLocal        []BatchFermentablesLocal `json:"batchFermentablesLocal"`
	MeasuredMashEfficiency        float64                  `json:"measuredMashEfficiency"`
	MeasuredKettleSize            float64                  `json:"measuredKettleSize"`
	BatchFermentables             []BatchFermentables      `json:"batchFermentables"`
	EstimatedBuGuRatio            float64                  `json:"estimatedBuGuRatio"`
	BatchYeasts                   []BatchYeasts            `json:"batchYeasts"`
	Version                       string                   `json:"_version"`
	BatchMiscsLocal               []interface{}            `json:"batchMiscsLocal"`
	EstimatedTotalGravity         float64                  `json:"estimatedTotalGravity"`
	EstimatedIbu                  int                      `json:"estimatedIbu"`
	BoilStepsCount                int                      `json:"boilStepsCount"`
	BatchHops                     []BatchHops              `json:"batchHops"`
	Rev                           string                   `json:"_rev"`
	FermentationStartDate         int64                    `json:"fermentationStartDate"`
	Hidden                        bool                     `json:"hidden"`
	Recipe                        Recipe                   `json:"recipe"`
	MeasuredOgSet                 bool                     `json:"measuredOgSet"`
	ID                            string                   `json:"_id"`
	Created                       Created                  `json:"_created"`
	FermentationControllerEnabled bool                     `json:"fermentationControllerEnabled"`
	PrimingSugarEquiv             interface{}              `json:"primingSugarEquiv"`
	Type                          string                   `json:"_type"`
	BatchNo                       int                      `json:"batchNo"`
	MeasuredPostBoilGravity       float64                  `json:"measuredPostBoilGravity"`
	BatchMiscs                    []BatchMiscs             `json:"batchMiscs"`
	MashStepsCount                int                      `json:"mashStepsCount"`
	TasteRating                   interface{}              `json:"tasteRating"`
	Events                        []Events                 `json:"events"`
	MeasuredAbv                   float64                  `json:"measuredAbv"`
	Brewer                        interface{}              `json:"brewer"`
	TimestampMs                   int64                    `json:"_timestamp_ms"`
	MeasuredPreBoilGravity        float64                  `json:"measuredPreBoilGravity"`
	BrewControllerEnabled         bool                     `json:"brewControllerEnabled"`
	EstimatedRbRatio              float64                  `json:"estimatedRbRatio"`
	EstimatedColor                float64                  `json:"estimatedColor"`
	MeasuredOg                    float64                  `json:"measuredOg"`
	MeasuredBoilSize              float64                  `json:"measuredBoilSize"`
	BatchYeastsLocal              []interface{}            `json:"batchYeastsLocal"`
	CarbonationForce              float64                  `json:"carbonationForce"`
	Cost                          Cost                     `json:"cost"`
	MeasuredConversionEfficiency  interface{}              `json:"measuredConversionEfficiency"`
	Measurements                  []interface{}            `json:"measurements"`
	Status                        string                   `json:"status"`
	Notes                         []Notes                  `json:"notes"`
	MeasuredAttenuation           float64                  `json:"measuredAttenuation"`
	EstimatedOg                   float64                  `json:"estimatedOg"`
	Init                          bool                     `json:"_init"`
	CarbonationType               string                   `json:"carbonationType"`
	BatchHopsLocal                []interface{}            `json:"batchHopsLocal"`
	HideBrewSheet                 bool                     `json:"hideBrewSheet"`
	BottlingDate                  int64                    `json:"bottlingDate"`
	BrewDate                      int64                    `json:"brewDate"`
	MeasuredKettleEfficiency      float64                  `json:"measuredKettleEfficiency"`
	EstimatedFg                   float64                  `json:"estimatedFg"`
	MeasuredMashPh                float64                  `json:"measuredMashPh"`
	Devices                       Devices                  `json:"devices"`
	MeasuredEfficiency            float64                  `json:"measuredEfficiency"`
	Timestamp                     Timestamp                `json:"_timestamp"`
	HideBatchEvents               bool                     `json:"hideBatchEvents"`
}

type BoilSteps struct {
	Time int    `json:"time"`
	Name string `json:"name"`
}

type BatchFermentablesLocal struct {
	PotentialPercentage interface{} `json:"potentialPercentage"`
	Use                 string      `json:"use"`
	GrainCategory       string      `json:"grainCategory"`
	NotInRecipe         bool        `json:"notInRecipe"`
	UserNotes           string      `json:"userNotes"`
	CostPerAmount       int         `json:"costPerAmount"`
	Amount              float64     `json:"amount"`
	Inventory           float64     `json:"inventory"`
	Name                string      `json:"name"`
	Attenuation         interface{} `json:"attenuation"`
	Supplier            string      `json:"supplier"`
	Type                string      `json:"type"`
	Percentage          float64     `json:"percentage"`
	Color               float64     `json:"color"`
	Origin              string      `json:"origin"`
	Potential           interface{} `json:"potential"`
	DisplayAmount       float64     `json:"displayAmount"`
	DiastaticPower      int         `json:"diastaticPower,omitempty"`
	ID                  string      `json:"_id"`
	Lovibond            float64     `json:"lovibond,omitempty"`
}

type Timestamp struct {
	Seconds     int `json:"_seconds"`
	Nanoseconds int `json:"_nanoseconds"`
}

type Created struct {
	Seconds     int `json:"_seconds"`
	Nanoseconds int `json:"_nanoseconds"`
}

type BatchFermentables struct {
	PotentialPercentage float64     `json:"potentialPercentage"`
	Name                string      `json:"name"`
	ID                  string      `json:"_id"`
	Attenuation         interface{} `json:"attenuation"`
	CostPerAmount       int         `json:"costPerAmount"`
	Type                string      `json:"type"`
	NotFermentable      interface{} `json:"notFermentable"`
	Amount              float64     `json:"amount"`
	TotalCost           int         `json:"totalCost"`
	Potential           float64     `json:"potential"`
	Inventory           float64     `json:"inventory"`
	Supplier            string      `json:"supplier"`
	Color               float64     `json:"color"`
	IbuPerAmount        interface{} `json:"ibuPerAmount"`
	Origin              string      `json:"origin"`
	DisplayAmount       float64     `json:"displayAmount"`
	NotInRecipe         bool        `json:"notInRecipe"`
	UsedIn              string      `json:"usedIn,omitempty"`
	Version             string      `json:"_version,omitempty"`
	Hidden              bool        `json:"hidden,omitempty"`
	Timestamp           Timestamp   `json:"_timestamp,omitempty"`
	ManufacturingDate   interface{} `json:"manufacturingDate,omitempty"`
	Moisture            interface{} `json:"moisture,omitempty"`
	MaxInBatch          int         `json:"maxInBatch,omitempty"`
	CoarseFineDiff      float64     `json:"coarseFineDiff,omitempty"`
	Created             Created     `json:"_created,omitempty"`
	Substitutes         string      `json:"substitutes,omitempty"`
	Rev                 string      `json:"_rev,omitempty"`
	Notes               string      `json:"notes,omitempty"`
	TimestampMs         int64       `json:"_timestamp_ms,omitempty"`
	Acid                int         `json:"acid,omitempty"`
	Protein             interface{} `json:"protein,omitempty"`
	DiastaticPower      interface{} `json:"diastaticPower,omitempty"`
	BestBeforeDate      interface{} `json:"bestBeforeDate,omitempty"`
	GrainCategory       string      `json:"grainCategory,omitempty"`
	UserNotes           string      `json:"userNotes,omitempty"`
}

type BatchYeasts struct {
	ID            string  `json:"_id"`
	NotInRecipe   bool    `json:"notInRecipe"`
	MaxTemp       float64 `json:"maxTemp"`
	FermentsAll   bool    `json:"fermentsAll"`
	Attenuation   int     `json:"attenuation"`
	TotalCost     int     `json:"totalCost"`
	InventoryUnit string  `json:"inventoryUnit"`
	Laboratory    string  `json:"laboratory"`
	Unit          string  `json:"unit"`
	Type          string  `json:"type"`
	CostPerAmount int     `json:"costPerAmount"`
	Form          string  `json:"form"`
	Flocculation  string  `json:"flocculation"`
	ProductID     string  `json:"productId"`
	MinTemp       float64 `json:"minTemp"`
	DisplayAmount float64 `json:"displayAmount"`
	Name          string  `json:"name"`
	Amount        float64 `json:"amount"`
	Inventory     float64 `json:"inventory"`
}

type BatchHops struct {
	Origin        string  `json:"origin"`
	NotInRecipe   bool    `json:"notInRecipe"`
	Alpha         float64 `json:"alpha"`
	CostPerAmount int     `json:"costPerAmount"`
	Type          string  `json:"type"`
	TotalCost     int     `json:"totalCost"`
	Inventory     float64 `json:"inventory"`
	Amount        float64 `json:"amount"`
	Usage         string  `json:"usage"`
	ID            string  `json:"_id"`
	Name          string  `json:"name"`
	DisplayAmount float64 `json:"displayAmount"`
}

type Miscs struct {
	Timestamp         Timestamp   `json:"_timestamp"`
	Created           Created     `json:"_created"`
	Concentration     interface{} `json:"concentration,omitempty"`
	TimeIsDays        bool        `json:"timeIsDays"`
	WaterAdjustment   bool        `json:"waterAdjustment"`
	Unit              string      `json:"unit"`
	Time              interface{} `json:"time"`
	ID                string      `json:"_id"`
	BestBeforeDate    interface{} `json:"bestBeforeDate"`
	ManufacturingDate interface{} `json:"manufacturingDate"`
	AmountPerL        interface{} `json:"amountPerL"`
	Inventory         float64     `json:"inventory"`
	UserNotes         string      `json:"userNotes"`
	Use               string      `json:"use"`
	Version           string      `json:"_version"`
	Substitutes       string      `json:"substitutes"`
	Rev               string      `json:"_rev"`
	Hidden            bool        `json:"hidden"`
	TimestampMs       int64       `json:"_timestamp_ms"`
	Amount            float64     `json:"amount"`
	Name              string      `json:"name"`
	Notes             string      `json:"notes"`
	Type              string      `json:"type"`
}

type Defaults struct {
	Weight      string `json:"weight"`
	Color       string `json:"color"`
	Temp        string `json:"temp"`
	Altitude    string `json:"altitude"`
	GrainColor  string `json:"grainColor"`
	Attenuation string `json:"attenuation"`
	Preferred   string `json:"preferred"`
	Abv         string `json:"abv"`
	Gravity     string `json:"gravity"`
	Carbonation string `json:"carbonation"`
	Volume      string `json:"volume"`
	Hop         string `json:"hop"`
	Ibu         string `json:"ibu"`
	Pressure    string `json:"pressure"`
}

type SpargeTargetDiff struct {
	Sulfate                    float64 `json:"sulfate"`
	IonBalanceOff              bool    `json:"ionBalanceOff"`
	Calcium                    float64 `json:"calcium"`
	Alkalinity                 float64 `json:"alkalinity"`
	IonBalance                 int     `json:"ionBalance"`
	Sodium                     int     `json:"sodium"`
	BicarbonateMeqL            float64 `json:"bicarbonateMeqL"`
	Hardness                   int     `json:"hardness"`
	Anions                     float64 `json:"anions"`
	Magnesium                  float64 `json:"magnesium"`
	ResidualAlkalinity         float64 `json:"residualAlkalinity"`
	ResidualAlkalinityMeqLCalc float64 `json:"residualAlkalinityMeqLCalc"`
	Chloride                   float64 `json:"chloride"`
	SoClRatio                  float64 `json:"soClRatio"`
	Bicarbonate                int     `json:"bicarbonate"`
	Cations                    float64 `json:"cations"`
}

type WaterMash struct {
	Created                    Created   `json:"_created"`
	IonBalance                 int       `json:"ionBalance"`
	Hidden                     bool      `json:"hidden"`
	Ph                         float64   `json:"ph"`
	Alkalinity                 float64   `json:"alkalinity"`
	Rev                        string    `json:"_rev"`
	Type                       string    `json:"type"`
	Cations                    float64   `json:"cations"`
	Name                       string    `json:"name"`
	ResidualAlkalinityMeqLCalc float64   `json:"residualAlkalinityMeqLCalc"`
	Anions                     float64   `json:"anions"`
	ID                         string    `json:"_id"`
	Sulfate                    float64   `json:"sulfate"`
	ResidualAlkalinity         float64   `json:"residualAlkalinity"`
	Chloride                   float64   `json:"chloride"`
	TimestampMs                int64     `json:"_timestamp_ms"`
	SoClRatio                  float64   `json:"soClRatio"`
	Version                    string    `json:"_version"`
	Bicarbonate                int       `json:"bicarbonate"`
	Magnesium                  float64   `json:"magnesium"`
	IonBalanceOff              bool      `json:"ionBalanceOff"`
	Timestamp                  Timestamp `json:"_timestamp"`
	Hardness                   int       `json:"hardness"`
	Calcium                    float64   `json:"calcium"`
	BicarbonateMeqL            float64   `json:"bicarbonateMeqL"`
	Sodium                     int       `json:"sodium"`
}

type MashTargetDiff struct {
	ResidualAlkalinityMeqLCalc float64 `json:"residualAlkalinityMeqLCalc"`
	Sodium                     int     `json:"sodium"`
	IonBalanceOff              bool    `json:"ionBalanceOff"`
	Cations                    float64 `json:"cations"`
	Magnesium                  float64 `json:"magnesium"`
	BicarbonateMeqL            float64 `json:"bicarbonateMeqL"`
	Hardness                   int     `json:"hardness"`
	ResidualAlkalinity         float64 `json:"residualAlkalinity"`
	Chloride                   float64 `json:"chloride"`
	Alkalinity                 float64 `json:"alkalinity"`
	SoClRatio                  float64 `json:"soClRatio"`
	Sulfate                    float64 `json:"sulfate"`
	IonBalance                 int     `json:"ionBalance"`
	Calcium                    float64 `json:"calcium"`
	Anions                     float64 `json:"anions"`
	Bicarbonate                int     `json:"bicarbonate"`
}

type TotalTargetDiff struct {
	ResidualAlkalinity         float64 `json:"residualAlkalinity"`
	Anions                     float64 `json:"anions"`
	Calcium                    float64 `json:"calcium"`
	ResidualAlkalinityMeqLCalc float64 `json:"residualAlkalinityMeqLCalc"`
	Hardness                   int     `json:"hardness"`
	Sodium                     int     `json:"sodium"`
	Chloride                   float64 `json:"chloride"`
	IonBalanceOff              bool    `json:"ionBalanceOff"`
	IonBalance                 int     `json:"ionBalance"`
	Cations                    float64 `json:"cations"`
	Magnesium                  float64 `json:"magnesium"`
	Bicarbonate                int     `json:"bicarbonate"`
	SoClRatio                  float64 `json:"soClRatio"`
	Alkalinity                 float64 `json:"alkalinity"`
	BicarbonateMeqL            float64 `json:"bicarbonateMeqL"`
	Sulfate                    float64 `json:"sulfate"`
}

type Source struct {
	Magnesium                  int       `json:"magnesium"`
	Created                    Created   `json:"_created"`
	Rev                        string    `json:"_rev"`
	ID                         string    `json:"_id"`
	Cations                    float64   `json:"cations"`
	Version                    string    `json:"_version"`
	Type                       string    `json:"type"`
	Hidden                     bool      `json:"hidden"`
	Hardness                   int       `json:"hardness"`
	Anions                     float64   `json:"anions"`
	Sulfate                    int       `json:"sulfate"`
	Sodium                     int       `json:"sodium"`
	Ph                         float64   `json:"ph"`
	Bicarbonate                int       `json:"bicarbonate"`
	TimestampMs                int64     `json:"_timestamp_ms"`
	IonBalanceOff              bool      `json:"ionBalanceOff"`
	Alkalinity                 float64   `json:"alkalinity"`
	ResidualAlkalinityMeqLCalc float64   `json:"residualAlkalinityMeqLCalc"`
	BicarbonateMeqL            float64   `json:"bicarbonateMeqL"`
	Calcium                    int       `json:"calcium"`
	Chloride                   int       `json:"chloride"`
	SoClRatio                  float64   `json:"soClRatio"`
	Name                       string    `json:"name"`
	IonBalance                 int       `json:"ionBalance"`
	Timestamp                  Timestamp `json:"_timestamp"`
	ResidualAlkalinity         float64   `json:"residualAlkalinity"`
}

type SpargeAdjustmentsAcids struct {
	Amount        int    `json:"amount"`
	Type          string `json:"type"`
	Concentration int    `json:"concentration"`
}

type SpargeAdjustments struct {
	Acids                  []SpargeAdjustmentsAcids `json:"acids"`
	Calcium                float64                  `json:"calcium"`
	Chloride               float64                  `json:"chloride"`
	LtAMS                  int                      `json:"ltAMS"`
	Volume                 float64                  `json:"volume"`
	SodiumBicarbonate      int                      `json:"sodiumBicarbonate"`
	SodiumChloride         int                      `json:"sodiumChloride"`
	CalciumCarbonate       int                      `json:"calciumCarbonate"`
	LtDWB                  int                      `json:"ltDWB"`
	CalciumChloride        float64                  `json:"calciumChloride"`
	MagnesiumSulfate       float64                  `json:"magnesiumSulfate"`
	SodiumMetabisulfite    int                      `json:"sodiumMetabisulfite"`
	Sodium                 int                      `json:"sodium"`
	CalciumHydroxide       int                      `json:"calciumHydroxide"`
	MagnesiumChloride      int                      `json:"magnesiumChloride"`
	SodiumMetabisulfitePPM int                      `json:"sodiumMetabisulfitePPM"`
	CalciumSulfate         float64                  `json:"calciumSulfate"`
	Bicarbonate            int                      `json:"bicarbonate"`
	Magnesium              float64                  `json:"magnesium"`
	Sulfate                float64                  `json:"sulfate"`
}

type WaterMeta struct {
	EqualSourceTotal bool `json:"equalSourceTotal"`
}

type Total struct {
	Cations                    float64   `json:"cations"`
	Alkalinity                 float64   `json:"alkalinity"`
	Anions                     float64   `json:"anions"`
	Calcium                    float64   `json:"calcium"`
	IonBalanceOff              bool      `json:"ionBalanceOff"`
	Sodium                     int       `json:"sodium"`
	Hardness                   int       `json:"hardness"`
	Bicarbonate                int       `json:"bicarbonate"`
	Version                    string    `json:"_version"`
	Rev                        string    `json:"_rev"`
	Name                       string    `json:"name"`
	Sulfate                    float64   `json:"sulfate"`
	Type                       string    `json:"type"`
	IonBalance                 int       `json:"ionBalance"`
	ResidualAlkalinity         float64   `json:"residualAlkalinity"`
	ID                         string    `json:"_id"`
	Chloride                   float64   `json:"chloride"`
	SoClRatio                  float64   `json:"soClRatio"`
	ResidualAlkalinityMeqLCalc float64   `json:"residualAlkalinityMeqLCalc"`
	BicarbonateMeqL            float64   `json:"bicarbonateMeqL"`
	Magnesium                  float64   `json:"magnesium"`
	Hidden                     bool      `json:"hidden"`
	Created                    Created   `json:"_created"`
	Ph                         float64   `json:"ph"`
	TimestampMs                int64     `json:"_timestamp_ms"`
	Timestamp                  Timestamp `json:"_timestamp"`
}

type MashAdjustmentsAcids struct {
	Type           string  `json:"type"`
	AlkalinityMeqL float64 `json:"alkalinityMeqL"`
	Concentration  int     `json:"concentration"`
	Amount         int     `json:"amount"`
}

type MashAdjustments struct {
	SodiumMetabisulfite    int                    `json:"sodiumMetabisulfite"`
	MagnesiumChloride      int                    `json:"magnesiumChloride"`
	Volume                 float64                `json:"volume"`
	Sodium                 int                    `json:"sodium"`
	SodiumBicarbonate      int                    `json:"sodiumBicarbonate"`
	MagnesiumSulfate       float64                `json:"magnesiumSulfate"`
	Magnesium              float64                `json:"magnesium"`
	Sulfate                float64                `json:"sulfate"`
	CalciumChloride        float64                `json:"calciumChloride"`
	CalciumHydroxide       int                    `json:"calciumHydroxide"`
	SodiumChloride         int                    `json:"sodiumChloride"`
	Bicarbonate            int                    `json:"bicarbonate"`
	Acids                  []MashAdjustmentsAcids `json:"acids"`
	LtDWB                  int                    `json:"ltDWB"`
	CalciumCarbonate       int                    `json:"calciumCarbonate"`
	CalciumSulfate         float64                `json:"calciumSulfate"`
	Chloride               float64                `json:"chloride"`
	SodiumMetabisulfitePPM int                    `json:"sodiumMetabisulfitePPM"`
	LtAMS                  int                    `json:"ltAMS"`
	Calcium                float64                `json:"calcium"`
}

type Target struct {
	Type                       string  `json:"type"`
	Bicarbonate                int     `json:"bicarbonate"`
	ResidualAlkalinity         float64 `json:"residualAlkalinity"`
	Alkalinity                 float64 `json:"alkalinity"`
	IonBalanceOff              bool    `json:"ionBalanceOff"`
	Magnesium                  int     `json:"magnesium"`
	Chloride                   int     `json:"chloride"`
	Sodium                     int     `json:"sodium"`
	Sulfate                    int     `json:"sulfate"`
	Anions                     float64 `json:"anions"`
	IonBalance                 int     `json:"ionBalance"`
	SoClRatio                  float64 `json:"soClRatio"`
	Cations                    float64 `json:"cations"`
	Calcium                    int     `json:"calcium"`
	BicarbonateMeqL            float64 `json:"bicarbonateMeqL"`
	ID                         string  `json:"_id"`
	Name                       string  `json:"name"`
	ResidualAlkalinityMeqLCalc float64 `json:"residualAlkalinityMeqLCalc"`
	Hardness                   int     `json:"hardness"`
}

type CalciumChloride struct {
	Sparge bool   `json:"sparge"`
	Auto   bool   `json:"auto"`
	Form   string `json:"form"`
	Mash   bool   `json:"mash"`
}

type CalciumHydroxide struct {
	Mash   bool `json:"mash"`
	Auto   bool `json:"auto"`
	Sparge bool `json:"sparge"`
}

type SodiumBicarbonate struct {
	Mash   bool `json:"mash"`
	Sparge bool `json:"sparge"`
	Auto   bool `json:"auto"`
}

type MagnesiumSulfate struct {
	Sparge bool `json:"sparge"`
	Auto   bool `json:"auto"`
	Mash   bool `json:"mash"`
}

type CalciumSulfate struct {
	Sparge bool `json:"sparge"`
	Mash   bool `json:"mash"`
	Auto   bool `json:"auto"`
}

type Settings struct {
	AdjustSparge      bool              `json:"adjustSparge"`
	CalciumChloride   CalciumChloride   `json:"calciumChloride"`
	CalciumHydroxide  CalciumHydroxide  `json:"calciumHydroxide"`
	SodiumBicarbonate SodiumBicarbonate `json:"sodiumBicarbonate"`
	MagnesiumSulfate  MagnesiumSulfate  `json:"magnesiumSulfate"`
	CalciumSulfate    CalciumSulfate    `json:"calciumSulfate"`
}

type TotalAdjustments struct {
	Calcium                float64 `json:"calcium"`
	Sulfate                float64 `json:"sulfate"`
	Bicarbonate            int     `json:"bicarbonate"`
	MagnesiumChloride      int     `json:"magnesiumChloride"`
	SodiumChloride         int     `json:"sodiumChloride"`
	LtAMS                  int     `json:"ltAMS"`
	CalciumCarbonate       int     `json:"calciumCarbonate"`
	LtDWB                  int     `json:"ltDWB"`
	SodiumBicarbonate      int     `json:"sodiumBicarbonate"`
	Sodium                 int     `json:"sodium"`
	Chloride               float64 `json:"chloride"`
	CalciumSulfate         float64 `json:"calciumSulfate"`
	SodiumMetabisulfite    int     `json:"sodiumMetabisulfite"`
	CalciumHydroxide       int     `json:"calciumHydroxide"`
	MagnesiumSulfate       float64 `json:"magnesiumSulfate"`
	SodiumMetabisulfitePPM int     `json:"sodiumMetabisulfitePPM"`
	CalciumChloride        float64 `json:"calciumChloride"`
	Volume                 float64 `json:"volume"`
	Magnesium              float64 `json:"magnesium"`
}

type SourceTargetDiff struct {
	IonBalance                 int     `json:"ionBalance"`
	Alkalinity                 float64 `json:"alkalinity"`
	Chloride                   int     `json:"chloride"`
	Sodium                     int     `json:"sodium"`
	Anions                     float64 `json:"anions"`
	Bicarbonate                int     `json:"bicarbonate"`
	Hardness                   int     `json:"hardness"`
	SoClRatio                  float64 `json:"soClRatio"`
	Sulfate                    int     `json:"sulfate"`
	BicarbonateMeqL            float64 `json:"bicarbonateMeqL"`
	IonBalanceOff              bool    `json:"ionBalanceOff"`
	Magnesium                  int     `json:"magnesium"`
	ResidualAlkalinityMeqLCalc float64 `json:"residualAlkalinityMeqLCalc"`
	Cations                    float64 `json:"cations"`
	ResidualAlkalinity         float64 `json:"residualAlkalinity"`
	Calcium                    int     `json:"calcium"`
}

type Sparge struct {
	Bicarbonate                int       `json:"bicarbonate"`
	Hardness                   int       `json:"hardness"`
	Created                    Created   `json:"_created"`
	BicarbonateMeqL            float64   `json:"bicarbonateMeqL"`
	Type                       string    `json:"type"`
	TimestampMs                int64     `json:"_timestamp_ms"`
	Alkalinity                 float64   `json:"alkalinity"`
	Hidden                     bool      `json:"hidden"`
	Anions                     float64   `json:"anions"`
	Magnesium                  float64   `json:"magnesium"`
	Name                       string    `json:"name"`
	Cations                    float64   `json:"cations"`
	Sulfate                    float64   `json:"sulfate"`
	Sodium                     int       `json:"sodium"`
	IonBalance                 int       `json:"ionBalance"`
	Version                    string    `json:"_version"`
	Rev                        string    `json:"_rev"`
	IonBalanceOff              bool      `json:"ionBalanceOff"`
	ResidualAlkalinity         float64   `json:"residualAlkalinity"`
	ID                         string    `json:"_id"`
	Chloride                   float64   `json:"chloride"`
	Calcium                    float64   `json:"calcium"`
	ResidualAlkalinityMeqLCalc float64   `json:"residualAlkalinityMeqLCalc"`
	Ph                         float64   `json:"ph"`
	Timestamp                  Timestamp `json:"_timestamp"`
	SoClRatio                  float64   `json:"soClRatio"`
}

type Water struct {
	SpargeTargetDiff        SpargeTargetDiff  `json:"spargeTargetDiff"`
	Mash                    WaterMash         `json:"mash"`
	MashTargetDiff          MashTargetDiff    `json:"mashTargetDiff"`
	TotalTargetDiff         TotalTargetDiff   `json:"totalTargetDiff"`
	EnableAcidAdjustments   bool              `json:"enableAcidAdjustments"`
	Source                  Source            `json:"source"`
	SpargeAdjustments       SpargeAdjustments `json:"spargeAdjustments"`
	MashPh                  float64           `json:"mashPh"`
	Meta                    WaterMeta         `json:"meta"`
	MashPhDistilled         float64           `json:"mashPhDistilled"`
	Total                   Total             `json:"total"`
	MashAdjustments         MashAdjustments   `json:"mashAdjustments"`
	MashWaterAmount         interface{}       `json:"mashWaterAmount"`
	AcidPhAdjustment        float64           `json:"acidPhAdjustment"`
	SpargeAcidPhAdjustment  int               `json:"spargeAcidPhAdjustment"`
	Diluted                 interface{}       `json:"diluted"`
	DilutionPercentage      interface{}       `json:"dilutionPercentage"`
	EnableSpargeAdjustments bool              `json:"enableSpargeAdjustments"`
	Target                  Target            `json:"target"`
	Style                   string            `json:"style"`
	Settings                Settings          `json:"settings"`
	SpargeWaterAmount       interface{}       `json:"spargeWaterAmount"`
	TotalAdjustments        TotalAdjustments  `json:"totalAdjustments"`
	SourceTargetDiff        SourceTargetDiff  `json:"sourceTargetDiff"`
	Sparge                  Sparge            `json:"sparge"`
}

type EquipmentMeta struct {
	MashEfficiencyIsCalculated bool `json:"mashEfficiencyIsCalculated"`
	EfficiencyIsCalculated     bool `json:"efficiencyIsCalculated"`
}

type Equipment struct {
	EfficiencyType              string        `json:"efficiencyType"`
	AromaHopUtilization         float64       `json:"aromaHopUtilization"`
	Version                     string        `json:"_version"`
	Meta                        EquipmentMeta `json:"_meta"`
	AltitudeAdjustment          bool          `json:"altitudeAdjustment"`
	WaterCalculation            string        `json:"waterCalculation"`
	MashWaterFormula            string        `json:"mashWaterFormula"`
	SpargeWaterReminderEnabled  bool          `json:"spargeWaterReminderEnabled"`
	FermenterVolume             float64       `json:"fermenterVolume"`
	Efficiency                  float64       `json:"efficiency"`
	BottlingVolume              float64       `json:"bottlingVolume"`
	FermenterTopUp              interface{}   `json:"fermenterTopUp"`
	Name                        string        `json:"name"`
	MashTunLoss                 float64       `json:"mashTunLoss"`
	HopstandTemperature         int           `json:"hopstandTemperature"`
	Timestamp                   Timestamp     `json:"_timestamp"`
	MashEfficiency              float64       `json:"mashEfficiency"`
	BoilTemp                    float64       `json:"boilTemp"`
	SpargeWaterFormula          string        `json:"spargeWaterFormula"`
	FermenterVolumeBeforeTopUp  float64       `json:"fermenterVolumeBeforeTopUp"`
	BoilExpansion               float64       `json:"boilExpansion"`
	GrainTemperature            float64       `json:"grainTemperature"`
	BrewhouseEfficiency         float64       `json:"brewhouseEfficiency"`
	FermenterLoss               float64       `json:"fermenterLoss"`
	MashTunDeadSpace            float64       `json:"mashTunDeadSpace"`
	BoilTime                    int           `json:"boilTime"`
	SpargeWaterOverflow         string        `json:"spargeWaterOverflow"`
	CalcBoilVolume              bool          `json:"calcBoilVolume"`
	ID                          string        `json:"_id"`
	EvaporationRate             float64       `json:"evaporationRate"`
	PostBoilKettleVol           float64       `json:"postBoilKettleVol"`
	MashTunHeatCapacity         float64       `json:"mashTunHeatCapacity"`
	HopUtilization              int           `json:"hopUtilization"`
	TimestampMs                 int64         `json:"_timestamp_ms"`
	CalcMashEfficiency          bool          `json:"calcMashEfficiency"`
	CalcAromaHopUtilization     bool          `json:"calcAromaHopUtilization"`
	Altitude                    float64       `json:"altitude"`
	Hidden                      bool          `json:"hidden"`
	AmbientTemperature          float64       `json:"ambientTemperature"`
	SpargeTemperature           float64       `json:"spargeTemperature"`
	Created                     Created       `json:"_created"`
	CalcStrikeWaterTemperature  bool          `json:"calcStrikeWaterTemperature"`
	MashWaterVolumeLimitEnabled bool          `json:"mashWaterVolumeLimitEnabled"`
	FermenterLossEstimate       int           `json:"fermenterLossEstimate"`
	BatchSize                   float64       `json:"batchSize"`
	Rev                         string        `json:"_rev"`
	BoilOffPerHr                float64       `json:"boilOffPerHr"`
	TrubChillerLoss             float64       `json:"trubChillerLoss"`
	BoilSize                    float64       `json:"boilSize"`
}

type Yeasts struct {
	Starter            bool        `json:"starter"`
	MaxTemp            float64     `json:"maxTemp"`
	Laboratory         string      `json:"laboratory"`
	MinAttenuation     interface{} `json:"minAttenuation"`
	Description        string      `json:"description"`
	Form               string      `json:"form"`
	Amount             int         `json:"float64"`
	Name               string      `json:"name"`
	Attenuation        int         `json:"attenuation"`
	FermentsAll        bool        `json:"fermentsAll"`
	BestBeforeDate     interface{} `json:"bestBeforeDate"`
	UserNotes          string      `json:"userNotes"`
	Type               string      `json:"type"`
	MinTemp            float64     `json:"minTemp"`
	Flocculation       string      `json:"flocculation"`
	ID                 string      `json:"_id"`
	ProductID          string      `json:"productId"`
	StarterGramExtract interface{} `json:"starterGramExtract"`
	ManufacturingDate  int64       `json:"manufacturingDate"`
	Unit               string      `json:"unit"`
	MaxAbv             interface{} `json:"maxAbv"`
	StarterSize        interface{} `json:"starterSize"`
	MaxAttenuation     interface{} `json:"maxAttenuation"`
}

type Hops struct {
	Type              string      `json:"type"`
	Time              int         `json:"time"`
	Usage             string      `json:"usage"`
	UsedIn            string      `json:"usedIn"`
	Humulene          interface{} `json:"humulene"`
	Alpha             float64     `json:"alpha"`
	Substitutes       string      `json:"substitutes"`
	Farnesene         interface{} `json:"farnesene"`
	BestBeforeDate    interface{} `json:"bestBeforeDate"`
	Hsi               interface{} `json:"hsi"`
	Temp              interface{} `json:"temp"`
	Use               string      `json:"use"`
	UserNotes         string      `json:"userNotes"`
	Notes             string      `json:"notes"`
	Name              string      `json:"name"`
	Beta              interface{} `json:"beta"`
	ManufacturingDate interface{} `json:"manufacturingDate"`
	Caryophyllene     interface{} `json:"caryophyllene"`
	Myrcene           interface{} `json:"myrcene"`
	Origin            string      `json:"origin"`
	Ibu               float64     `json:"ibu"`
	Oil               interface{} `json:"oil"`
	Year              interface{} `json:"year"`
	Amount            float64     `json:"amount"`
	ID                string      `json:"_id"`
	Cohumulone        interface{} `json:"cohumulone"`
}

type YeastSteps struct {
	HarvestRatio             int     `json:"harvestRatio"`
	StarterNewCellCount      float64 `json:"starterNewCellCount"`
	StarterPitchCellCount    float64 `json:"starterPitchCellCount"`
	GrowthFactor             float64 `json:"growthFactor"`
	StarterMillionCellsPerMl int     `json:"starterMillionCellsPerMl"`
	StarterGramLiquidExtract int     `json:"starterGramLiquidExtract"`
	StartGravity             float64 `json:"startGravity"`
	StarterGramExtract       int     `json:"starterGramExtract"`
	StarterPitchRate         float64 `json:"starterPitchRate"`
	StartCells               int     `json:"startCells"`
	StartVol                 float64 `json:"startVol"`
	InoculationRate          float64 `json:"inoculationRate"`
}

type Yeast struct {
	InitCells         int          `json:"initCells"`
	Mode              string       `json:"mode"`
	MillionCellsPerMl int          `json:"millionCellsPerMl"`
	ManufacturingDate int64        `json:"manufacturingDate"`
	PitchRatePkg      float64      `json:"pitchRatePkg"`
	CellsPrGramDry    int          `json:"cellsPrGramDry"`
	Og                float64      `json:"og"`
	StartType         string       `json:"startType"`
	OverbuildCells    interface{}  `json:"overbuildCells"`
	StarterPackages   float64      `json:"starterPackages"`
	CalcStarterVol    bool         `json:"calcStarterVol"`
	PacksPitchCells   int          `json:"packsPitchCells"`
	Rehydrate         bool         `json:"rehydrate"`
	Type              string       `json:"type"`
	CellsPerLiquidPkg int          `json:"cellsPerLiquidPkg"`
	Viability         float64      `json:"viability"`
	GramsPerDryPkg    float64      `json:"gramsPerDryPkg"`
	Packs             int          `json:"packs"`
	Volume            float64      `json:"volume"`
	CalcViablity      bool         `json:"calcViablity"`
	GrowthModel       string       `json:"growthModel"`
	PurePitch         bool         `json:"purePitch"`
	PitchRateActual   float64      `json:"pitchRateActual"`
	ActualPacks       float64      `json:"actualPacks"`
	PitchCells        float64      `json:"pitchCells"`
	Rate              float64      `json:"rate"`
	Steps             []YeastSteps `json:"steps"`
}

type Style struct {
	RbrMax           float64 `json:"rbrMax"`
	CarbonationStyle string  `json:"carbonationStyle"`
	IbuMax           int     `json:"ibuMax"`
	ColorMax         int     `json:"colorMax"`
	StyleLetter      string  `json:"styleLetter"`
	CategoryNumber   string  `json:"categoryNumber"`
	FgMin            float64 `json:"fgMin"`
	LovibondMax      int     `json:"lovibondMax"`
	ColorMin         int     `json:"colorMin"`
	AbvMax           float64 `json:"abvMax"`
	FgMax            float64 `json:"fgMax"`
	Category         string  `json:"category"`
	BuGuMin          float64 `json:"buGuMin"`
	StyleGuide       string  `json:"styleGuide"`
	Name             string  `json:"name"`
	Type             string  `json:"type"`
	OgMax            float64 `json:"ogMax"`
	ID               string  `json:"_id"`
	IbuMin           int     `json:"ibuMin"`
	LovibondMin      int     `json:"lovibondMin"`
	RbrMin           float64 `json:"rbrMin"`
	AbvMin           float64 `json:"abvMin"`
	BuGuMax          float64 `json:"buGuMax"`
	OgMin            float64 `json:"ogMin"`
}

type OtherFermentables struct {
	Cgdb                interface{} `json:"cgdb"`
	UserNotes           string      `json:"userNotes"`
	Attenuation         float64     `json:"attenuation"`
	MaxInBatch          interface{} `json:"maxInBatch"`
	Fan                 interface{} `json:"fan"`
	Type                string      `json:"type"`
	Supplier            string      `json:"supplier"`
	Percentage          float64     `json:"percentage"`
	Potential           float64     `json:"potential"`
	DiastaticPower      interface{} `json:"diastaticPower"`
	GrainCategory       interface{} `json:"grainCategory"`
	Acid                interface{} `json:"acid"`
	Protein             interface{} `json:"protein"`
	CoarseFineDiff      interface{} `json:"coarseFineDiff"`
	ManufacturingDate   interface{} `json:"manufacturingDate"`
	CostPerAmount       interface{} `json:"costPerAmount"`
	Name                string      `json:"name"`
	Friability          interface{} `json:"friability"`
	Color               float64     `json:"color"`
	Moisture            interface{} `json:"moisture"`
	Amount              float64     `json:"amount"`
	Inventory           interface{} `json:"inventory"`
	Origin              string      `json:"origin"`
	ID                  string      `json:"_id"`
	Substitutes         string      `json:"substitutes"`
	Notes               string      `json:"notes"`
	UsedIn              string      `json:"usedIn"`
	Fgdb                interface{} `json:"fgdb"`
	NotFermentable      bool        `json:"notFermentable"`
	PotentialPercentage float64     `json:"potentialPercentage"`
	IbuPerAmount        interface{} `json:"ibuPerAmount"`
	BestBeforeDate      interface{} `json:"bestBeforeDate"`
}

type MashFermentables struct {
	Color               float64     `json:"color"`
	Supplier            string      `json:"supplier"`
	Percentage          float64     `json:"percentage"`
	UserNotes           string      `json:"userNotes"`
	Use                 string      `json:"use,omitempty"`
	Potential           interface{} `json:"potential"`
	Attenuation         interface{} `json:"attenuation"`
	DiastaticPower      int         `json:"diastaticPower,omitempty"`
	Amount              float64     `json:"amount"`
	Name                string      `json:"name"`
	PotentialPercentage interface{} `json:"potentialPercentage"`
	ID                  string      `json:"_id"`
	GrainCategory       string      `json:"grainCategory"`
	Type                string      `json:"type"`
	Origin              string      `json:"origin"`
	Lovibond            float64     `json:"lovibond,omitempty"`
	NotFermentable      bool        `json:"notFermentable,omitempty"`
	Friability          interface{} `json:"friability,omitempty"`
	Rev                 string      `json:"_rev,omitempty"`
	Timestamp           Timestamp   `json:"_timestamp,omitempty"`
	Version             string      `json:"_version,omitempty"`
	TimestampMs         int64       `json:"_timestamp_ms,omitempty"`
	ManufacturingDate   interface{} `json:"manufacturingDate,omitempty"`
	Created             Created     `json:"_created,omitempty"`
	Fan                 interface{} `json:"fan,omitempty"`
	BestBeforeDate      interface{} `json:"bestBeforeDate,omitempty"`
	Moisture            interface{} `json:"moisture,omitempty"`
	Inventory           float64     `json:"inventory,omitempty"`
	Hidden              bool        `json:"hidden,omitempty"`
	UsedIn              string      `json:"usedIn,omitempty"`
	IbuPerAmount        interface{} `json:"ibuPerAmount,omitempty"`
	Protein             interface{} `json:"protein,omitempty"`
	CostPerAmount       interface{} `json:"costPerAmount,omitempty"`
	Acid                int         `json:"acid,omitempty"`
	CoarseFineDiff      float64     `json:"coarseFineDiff,omitempty"`
	Cgdb                interface{} `json:"cgdb,omitempty"`
	Notes               string      `json:"notes,omitempty"`
	Fgdb                interface{} `json:"fgdb,omitempty"`
	MaxInBatch          int         `json:"maxInBatch,omitempty"`
	Substitutes         string      `json:"substitutes,omitempty"`
}

type Data struct {
	MashFermentablesAmount  float64             `json:"mashFermentablesAmount"`
	MashVolumeSurplus       int                 `json:"mashVolumeSurplus"`
	StrikeTemp              float64             `json:"strikeTemp"`
	AllDiastaticPower       bool                `json:"allDiastaticPower"`
	SpargeWaterAmount       float64             `json:"spargeWaterAmount"`
	BatchSpargeWaterAmount4 interface{}         `json:"batchSpargeWaterAmount4"`
	TopUpWater              int                 `json:"topUpWater"`
	OtherFermentables       []OtherFermentables `json:"otherFermentables"`
	MashFermentables        []MashFermentables  `json:"mashFermentables"`
	BatchSpargeWaterAmount1 interface{}         `json:"batchSpargeWaterAmount1"`
	HopsAmount              float64             `json:"hopsAmount"`
	OtherFermentablesAmount float64             `json:"otherFermentablesAmount"`
	MashVolume              float64             `json:"mashVolume"`
	MashWaterAmount         float64             `json:"mashWaterAmount"`
	HltWaterAmount          float64             `json:"hltWaterAmount"`
	TotalWaterAmount        float64             `json:"totalWaterAmount"`
	BatchSpargeWaterAmount3 interface{}         `json:"batchSpargeWaterAmount3"`
	BatchSpargeWaterAmount2 interface{}         `json:"batchSpargeWaterAmount2"`
	TotalDiastaticPower     float64             `json:"totalDiastaticPower"`
}

type Calories struct {
	Total   float64 `json:"total"`
	Carbs   float64 `json:"carbs"`
	KJ      float64 `json:"kJ"`
	Alcohol float64 `json:"alcohol"`
}

type Carbs struct {
	Total float64 `json:"total"`
}

type Nutrition struct {
	Calories Calories `json:"calories"`
	Carbs    Carbs    `json:"carbs"`
}

type MashSteps struct {
	DisplayStepTemp float64     `json:"displayStepTemp"`
	StepTime        int         `json:"stepTime"`
	RampTime        interface{} `json:"rampTime"`
	StepTemp        float64     `json:"stepTemp"`
	Name            string      `json:"name"`
	Type            string      `json:"type"`
}

type RecipeMash struct {
	Name  string      `json:"name"`
	Steps []MashSteps `json:"steps"`
	ID    string      `json:"_id"`
}

type FermentationSteps struct {
	StepTime        int         `json:"stepTime"`
	Type            string      `json:"type"`
	Ramp            interface{} `json:"ramp"`
	ActualTime      int64       `json:"actualTime"`
	DisplayStepTemp float64     `json:"displayStepTemp"`
	Pressure        interface{} `json:"pressure"`
	DisplayPressure interface{} `json:"displayPressure"`
	StepTemp        float64     `json:"stepTemp"`
}

type Fermentation struct {
	ID    string              `json:"_id"`
	Steps []FermentationSteps `json:"steps"`
	Name  string              `json:"name"`
}

type Fermentables struct {
	Attenuation         interface{} `json:"attenuation"`
	GrainCategory       string      `json:"grainCategory"`
	Type                string      `json:"type"`
	PotentialPercentage interface{} `json:"potentialPercentage"`
	UserNotes           string      `json:"userNotes"`
	Potential           interface{} `json:"potential"`
	Origin              string      `json:"origin"`
	DiastaticPower      int         `json:"diastaticPower,omitempty"`
	Supplier            string      `json:"supplier"`
	Percentage          float64     `json:"percentage"`
	ID                  string      `json:"_id"`
	Amount              float64     `json:"amount"`
	Name                string      `json:"name"`
	Use                 string      `json:"use,omitempty"`
	Color               float64     `json:"color"`
	Lovibond            float64     `json:"lovibond,omitempty"`
	BestBeforeDate      interface{} `json:"bestBeforeDate,omitempty"`
	Cgdb                interface{} `json:"cgdb,omitempty"`
	Protein             interface{} `json:"protein,omitempty"`
	Fgdb                interface{} `json:"fgdb,omitempty"`
	IbuPerAmount        interface{} `json:"ibuPerAmount,omitempty"`
	Moisture            interface{} `json:"moisture,omitempty"`
	NotFermentable      bool        `json:"notFermentable,omitempty"`
	Notes               string      `json:"notes,omitempty"`
	Fan                 interface{} `json:"fan,omitempty"`
	Friability          interface{} `json:"friability,omitempty"`
	Inventory           interface{} `json:"inventory,omitempty"`
	UsedIn              string      `json:"usedIn,omitempty"`
	ManufacturingDate   interface{} `json:"manufacturingDate,omitempty"`
	MaxInBatch          interface{} `json:"maxInBatch,omitempty"`
	CoarseFineDiff      interface{} `json:"coarseFineDiff,omitempty"`
	Acid                interface{} `json:"acid,omitempty"`
	CostPerAmount       interface{} `json:"costPerAmount,omitempty"`
	Substitutes         string      `json:"substitutes,omitempty"`
	Rev                 string      `json:"_rev,omitempty"`
	Created             Created     `json:"_created,omitempty"`
	Timestamp           Timestamp   `json:"_timestamp,omitempty"`
	TimestampMs         int64       `json:"_timestamp_ms,omitempty"`
	Version             string      `json:"_version,omitempty"`
	Hidden              bool        `json:"hidden,omitempty"`
}

type CarbonationStyle struct {
	CarbMax float64 `json:"carbMax"`
	CarbMin float64 `json:"carbMin"`
	ID      string  `json:"_id"`
	Name    string  `json:"name"`
}

type Recipe struct {
	ExtraGravity            int           `json:"extraGravity"`
	Version                 string        `json:"_version"`
	StyleAbv                bool          `json:"styleAbv"`
	Miscs                   []Miscs       `json:"miscs"`
	StyleColor              bool          `json:"styleColor"`
	MashEfficiency          float64       `json:"mashEfficiency"`
	Carbonation             float64       `json:"carbonation"`
	Og                      float64       `json:"og"`
	Tags                    interface{}   `json:"tags"`
	Defaults                Defaults      `json:"defaults"`
	Water                   Water         `json:"water"`
	RbRatio                 float64       `json:"rbRatio"`
	SumDryHopPerLiter       float64       `json:"sumDryHopPerLiter"`
	ID                      string        `json:"_id"`
	Ibu                     float64       `json:"ibu"`
	Attenuation             float64       `json:"attenuation"`
	Equipment               Equipment     `json:"equipment"`
	Hidden                  bool          `json:"hidden"`
	Yeasts                  []Yeasts      `json:"yeasts"`
	Hops                    []Hops        `json:"hops"`
	BoilSize                float64       `json:"boilSize"`
	Yeast                   Yeast         `json:"yeast"`
	PrimaryTemp             float64       `json:"primaryTemp"`
	Color                   float64       `json:"color"`
	Path                    string        `json:"path"`
	HopStandMinutes         int           `json:"hopStandMinutes"`
	BuGuRatio               float64       `json:"buGuRatio"`
	Author                  interface{}   `json:"author"`
	BatchSize               float64       `json:"batchSize"`
	DiastaticPower          int           `json:"diastaticPower"`
	TotalGravity            float64       `json:"totalGravity"`
	SearchTags              []interface{} `json:"searchTags"`
	Name                    string        `json:"name"`
	TimestampMs             int64         `json:"_timestamp_ms"`
	AvgWeightedHopstandTemp int           `json:"avgWeightedHopstandTemp"`
	StyleFg                 bool          `json:"styleFg"`
	Public                  bool          `json:"public"`
	FirstWortGravity        interface{}   `json:"firstWortGravity"`
	FgFormula               string        `json:"fgFormula"`
	Init                    bool          `json:"_init"`
	Ev                      float64       `json:"_ev"`
	HopsTotalAmount         float64       `json:"hopsTotalAmount"`
	// Type                     string           `json:"_type"` // TODO: Are these needed
	// Public                   bool             `json:"_public"`
	StyleIbu                 bool             `json:"styleIbu"`
	Style                    Style            `json:"style"`
	StyleConformity          bool             `json:"styleConformity"`
	IbuFormula               string           `json:"ibuFormula"`
	StyleBuGu                bool             `json:"styleBuGu"`
	StyleRbr                 bool             `json:"styleRbr"`
	Abv                      float64          `json:"abv"`
	Data                     Data             `json:"data"`
	Nutrition                Nutrition        `json:"nutrition"`
	Mash                     RecipeMash       `json:"mash"`
	PreBoilGravity           float64          `json:"preBoilGravity"`
	OgPlato                  float64          `json:"ogPlato"`
	FermentableIbu           int              `json:"fermentableIbu"`
	Fermentation             Fermentation     `json:"fermentation"`
	StyleCarb                bool             `json:"styleCarb"`
	Timestamp                Timestamp        `json:"_timestamp"`
	Type                     string           `json:"type"`
	FermentablesTotalAmount  float64          `json:"fermentablesTotalAmount"`
	Fermentables             []Fermentables   `json:"fermentables"`
	CarbonationStyle         CarbonationStyle `json:"carbonationStyle"`
	PostBoilGravity          float64          `json:"postBoilGravity"`
	Rev                      string           `json:"_rev"`
	StyleOg                  bool             `json:"styleOg"`
	BoilTime                 int              `json:"boilTime"`
	Fg                       float64          `json:"fg"`
	FgEstimated              float64          `json:"fgEstimated"`
	Created                  Created          `json:"_created"`
	YeastToleranceExceededBy interface{}      `json:"yeastToleranceExceededBy"`
	Efficiency               float64          `json:"efficiency"`
}

type BatchMiscs struct {
	TotalCost            int         `json:"totalCost"`
	Type                 string      `json:"type"`
	AmountPerL           interface{} `json:"amountPerL"`
	TimeIsDays           bool        `json:"timeIsDays"`
	Name                 string      `json:"name"`
	Rev                  string      `json:"_rev"`
	DisplayAmount        float64     `json:"displayAmount"`
	Created              Created     `json:"_created"`
	BestBeforeDate       interface{} `json:"bestBeforeDate"`
	CostPerAmount        int         `json:"costPerAmount"`
	Amount               float64     `json:"amount"`
	Hidden               bool        `json:"hidden"`
	WaterAdjustment      bool        `json:"waterAdjustment"`
	NotInRecipe          bool        `json:"notInRecipe"`
	Unit                 string      `json:"unit"`
	Inventory            float64     `json:"inventory"`
	Time                 interface{} `json:"time"`
	Notes                string      `json:"notes"`
	TimestampMs          int64       `json:"_timestamp_ms"`
	InventoryUnit        string      `json:"inventoryUnit"`
	Timestamp            Timestamp   `json:"_timestamp"`
	Substitutes          string      `json:"substitutes"`
	UserNotes            string      `json:"userNotes"`
	ID                   string      `json:"_id"`
	Use                  string      `json:"use"`
	ManufacturingDate    interface{} `json:"manufacturingDate"`
	Version              string      `json:"_version"`
	RemovedAmount        float64     `json:"removedAmount,omitempty"`
	Checked              bool        `json:"checked,omitempty"`
	RemovedFromInventory bool        `json:"removedFromInventory,omitempty"`
}

type Events struct {
	Active          bool   `json:"active"`
	NotifyTime      int64  `json:"notifyTime,omitempty"`
	DescriptionHTML string `json:"descriptionHTML"`
	DayEvent        bool   `json:"dayEvent"`
	EventText       string `json:"eventText"`
	EventType       string `json:"eventType"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Time            int64  `json:"time"`
}

type Cost struct {
	PerBottlingLiter  int `json:"perBottlingLiter"`
	Hops              int `json:"hops"`
	YeastsShare       int `json:"yeastsShare"`
	FermentablesShare int `json:"fermentablesShare"`
	Total             int `json:"total"`
	HopsShare         int `json:"hopsShare"`
	MiscsShare        int `json:"miscsShare"`
	Fermentables      int `json:"fermentables"`
	Miscs             int `json:"miscs"`
	Yeasts            int `json:"yeasts"`
}

type Notes struct {
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
	Type      string `json:"type,omitempty"`
	Note      string `json:"note"`
}

type MyBrewbot struct {
	Items   []interface{} `json:"items"`
	Enabled bool          `json:"enabled"`
}

type LastData struct {
	Type      string  `json:"type"`
	Timepoint int64   `json:"timepoint"`
	Time      int64   `json:"time"`
	Temp      float64 `json:"temp"`
	Status    string  `json:"status"`
	ID        string  `json:"id"`
	Sg        float64 `json:"sg"`
	Comment   string  `json:"comment"`
}

type Items struct {
	Key      string      `json:"key"`
	BatchID  interface{} `json:"batchId"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Enabled  bool        `json:"enabled"`
	Settings interface{} `json:"settings"`
	LastLog  int         `json:"lastLog"`
	LastData LastData    `json:"lastData"`
	Series   interface{} `json:"series"`
	Hidden   bool        `json:"hidden"`
}

type Tilt struct {
	Temp    bool    `json:"temp"`
	Mode    string  `json:"mode"`
	Gravity bool    `json:"gravity"`
	Enabled bool    `json:"enabled"`
	Items   []Items `json:"items"`
}

type SmartPid struct {
	Enabled      bool          `json:"enabled"`
	Items        []interface{} `json:"items"`
	BrewDeviceID interface{}   `json:"brewDeviceId"`
}

type FloatHydrometer struct {
	Enabled bool          `json:"enabled"`
	Items   []interface{} `json:"items"`
}

type ISpindel struct {
	Items   []interface{} `json:"items"`
	Enabled bool          `json:"enabled"`
}

type PlaatoKeg struct {
	Enabled bool          `json:"enabled"`
	Items   []interface{} `json:"items"`
}

type PlaatoAirlock struct {
	Items   []interface{} `json:"items"`
	Enabled bool          `json:"enabled"`
}

type Stream struct {
	Enabled bool          `json:"enabled"`
	Items   []interface{} `json:"items"`
}

type Gfcc struct {
	Items        []interface{} `json:"items"`
	Enabled      bool          `json:"enabled"`
	BrewDeviceID interface{}   `json:"brewDeviceId"`
}

type FloatyHydrometer struct {
	Enabled bool          `json:"enabled"`
	Items   []interface{} `json:"items"`
}

type BrewPiLess struct {
	Enabled bool          `json:"enabled"`
	Items   []interface{} `json:"items"`
}

type Devices struct {
	MyBrewbot        MyBrewbot        `json:"myBrewbot"`
	Tilt             Tilt             `json:"tilt"`
	SmartPid         SmartPid         `json:"smartPid"`
	FloatHydrometer  FloatHydrometer  `json:"floatHydrometer"`
	ISpindel         ISpindel         `json:"iSpindel"`
	PlaatoKeg        PlaatoKeg        `json:"plaatoKeg"`
	PlaatoAirlock    PlaatoAirlock    `json:"plaatoAirlock"`
	Stream           Stream           `json:"stream"`
	Gfcc             Gfcc             `json:"gfcc"`
	FloatyHydrometer FloatyHydrometer `json:"floatyHydrometer"`
	BrewPiLess       BrewPiLess       `json:"brewPiLess"`
}
