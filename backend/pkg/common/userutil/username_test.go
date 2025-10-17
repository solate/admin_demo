package userutil

import "testing"

func TestGenerateDefaultUsername(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected string
	}{
		{
			name:     "正常手机号",
			phone:    "13812348375",
			expected: "课研所用户8375",
		},
		{
			name:     "带前缀的手机号",
			phone:    "+8613812348375",
			expected: "课研所用户8375",
		},
		{
			name:     "11位手机号",
			phone:    "18888888888",
			expected: "课研所用户8888",
		},
		{
			name:     "空手机号",
			phone:    "",
			expected: "课研所用户",
		},
		{
			name:     "短手机号",
			phone:    "123",
			expected: "课研所用户",
		},
		{
			name:     "带空格的手机号",
			phone:    " 13812348375 ",
			expected: "课研所用户8375",
		},
		{
			name:     "恰好四位数字",
			phone:    "1234",
			expected: "课研所用户1234",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GenerateDefaultUsername(tt.phone)
			if result != tt.expected {
				t.Errorf("GenerateDefaultUsername(%q) = %q, want %q", tt.phone, result, tt.expected)
			}
		})
	}
}

func TestMaskPhone(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected string
	}{
		{
			name:     "正常11位手机号",
			phone:    "17612345678",
			expected: "176****5678",
		},
		{
			name:     "另一个11位手机号",
			phone:    "13812348375",
			expected: "138****8375",
		},
		{
			name:     "带空格的手机号",
			phone:    " 18888888888 ",
			expected: "188****8888",
		},
		{
			name:     "空手机号",
			phone:    "",
			expected: "",
		},
		{
			name:     "短于11位的号码",
			phone:    "123456789",
			expected: "123456789",
		},
		{
			name:     "长于11位的号码",
			phone:    "123456789012",
			expected: "123****9012",
		},
		{
			name:     "恰好11位",
			phone:    "12345678901",
			expected: "123****8901",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MaskPhone(tt.phone)
			if result != tt.expected {
				t.Errorf("MaskPhone(%q) = %q, want %q", tt.phone, result, tt.expected)
			}
		})
	}
}

func TestNormalizePhone(t *testing.T) {
	tests := []struct {
		name     string
		phone    string
		expected string
	}{
		{
			name:     "带+86-前缀的手机号",
			phone:    "+86-18270548941",
			expected: "18270548941",
		},
		{
			name:     "带+86 前缀的手机号（有空格）",
			phone:    "+86 18270548941",
			expected: "18270548941",
		},
		{
			name:     "带+86前缀的手机号（无分隔符）",
			phone:    "+8618270548941",
			expected: "18270548941",
		},
		{
			name:     "正常11位手机号",
			phone:    "18270548941",
			expected: "18270548941",
		},
		{
			name:     "带空格的手机号",
			phone:    " 18270548941 ",
			expected: "18270548941",
		},
		{
			name:     "带横线的手机号",
			phone:    "182-7054-8941",
			expected: "18270548941",
		},
		{
			name:     "混合格式的手机号",
			phone:    "+86-182 7054 8941",
			expected: "18270548941",
		},
		{
			name:     "空手机号",
			phone:    "",
			expected: "",
		},
		{
			name:     "只有空格的手机号",
			phone:    "   ",
			expected: "",
		},
		{
			name:     "包含字母的手机号",
			phone:    "+86-182abc70548941",
			expected: "18270548941",
		},
		{
			name:     "包含特殊字符的手机号",
			phone:    "+86-(182)-7054-8941",
			expected: "18270548941",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizePhone(tt.phone)
			if result != tt.expected {
				t.Errorf("NormalizePhone(%q) = %q, want %q", tt.phone, result, tt.expected)
			}
		})
	}
}
