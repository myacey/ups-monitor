package main

type UPSStatus struct {
	InputVoltage      float64
	InputFaultVoltage float64
	OutputVoltage     float64
	OutputCurrentPct  int
	InputFrequency    float64
	BatteryVoltage    float64
	Temperature       float64

	// Status bits
	UtilityFail  bool
	BatteryLow   bool
	BypassActive bool
	UPSFailed    bool
	IsStandby    bool
	TestRunning  bool
	Shutdown     bool
	BeeperOn     bool
}
