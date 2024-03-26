package main

import (
	"reflect"
	"testing"
)

func TestCollectTranslations(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "testing",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			translations, err := CollectTranslations("", []string{"en", "ru"}, nil)
			if err != nil {
				return
			}
			t.Log(translations)
		})
	}
}

func TestParseTemplatesDirectory(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "templates",
			args: args{dir: "/Volumes/WORK/BWG/iaac/gateways/admin/server/templates"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ExtractTemplatesDirectory(tt.args.dir, []string{"en", "ru"}, "_templ.go"); (err != nil) != tt.wantErr {
				t.Errorf("ExtractTemplatesDirectory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExtractLocalizationsDirectory(t *testing.T) {
	type args struct {
		dir string
	}
	tests := []struct {
		name    string
		args    args
		want    Language
		wantErr bool
	}{
		{
			name: "parse",
			args: args{dir: "/Volumes/WORK/BWG/iaac/gateways/admin/server/localizations"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractLocalizationsDirectory(tt.args.dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractLocalizationsDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtractLocalizationsDirectory() got = %v, want %v", got, tt.want)
			}
		})
	}
}
