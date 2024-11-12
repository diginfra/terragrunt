package options

// LevelFormatOptionName is the option name.
const LevelFormatOptionName = "format"

const (
	LevelFormatFull LevelFormatValue = iota
	LevelFormatShort
	LevelFormatTiny
)

var levelFormatList = NewMapValue(map[LevelFormatValue]string{ //nolint:gochecknoglobals
	LevelFormatTiny:  "tiny",
	LevelFormatShort: "short",
})

type LevelFormatValue byte

type LevelFormatOption struct {
	*CommonOption[LevelFormatValue]
}

// Format implements `Option` interface.
func (format *LevelFormatOption) Format(data *Data, _ string) (string, error) {
	switch format.value.Get() {
	case LevelFormatTiny:
		return data.Level.TinyName(), nil
	case LevelFormatShort:
		return data.Level.ShortName(), nil
	case LevelFormatFull:
	}

	return data.Level.FullName(), nil
}

// LevelFormat creates the option to format level name.
func LevelFormat(val LevelFormatValue) Option {
	return &LevelFormatOption{
		CommonOption: NewCommonOption(LevelFormatOptionName, levelFormatList.Set(val)),
	}
}
