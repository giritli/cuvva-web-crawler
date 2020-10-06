package sitemap

import (
	"net/url"
	"reflect"
	"testing"
)

func TestMap_Add(t *testing.T) {
	type args struct {
		u     *url.URL
		links []*url.URL
	}
	tests := []struct {
		name   string
		args   args
		ret    *Map
	}{
		{
			"generate nested map",
			args{
				u: &url.URL{
					Scheme:     "https",
					Host:       "test",
					Path:       "a/b/c",
				},
				links: nil,
			},
			&Map{
				Hosts: map[string]Page{
					"test": {
						Assets:   []string{},
						Children: map[string][]Page{
							"/a": {
								{
									Assets:   []string{},
									Children: map[string][]Page{
										"/b": {
											{
												Assets: []string{},
												Children: map[string][]Page{
													"/c": {
														{
															Assets: []string{},
															Children: map[string][]Page{},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"map with assets",
			args{
				u: &url.URL{
					Scheme:     "https",
					Host:       "test",
					Path:       "a",
				},
				links: []*url.URL{
					{
						Scheme:     "https",
						Host:       "test",
						Path:       "asset.jpg",
					},
				},
			},
			&Map{
				Hosts: map[string]Page{
					"test": {
						Assets: []string{},
						Children: map[string][]Page{
							"/a": {
								{
									Assets: []string{
										"https://test/asset.jpg",
									},
									Children: map[string][]Page{},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMap()
			m.Add(tt.args.u, tt.args.links)
			if !reflect.DeepEqual(m, tt.ret) {
				t.Errorf("Add() = %v, want %v", m, tt.ret)
			}
		})
	}
}
