package mac2vendor

func init() {
    {{- range $key, $value := . }}
    mapping["{{ $key }}"] = "{{ $value }}"
    {{- end }}
}