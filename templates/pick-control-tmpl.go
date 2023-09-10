package templates
var (
  PickQueryTemplate = `
{
  "__CONTENT": [

      {
          "id": {{.ShortQueryId}},
          "kind": "multiselect",
          "prompt": "enter ...",
            {{.MoveItemsInfo}}"
              {{- range $k, $v := . -}}
                {{ $v }}
                  "selector": [
                    {{.Index}},
                  ],
                  "children": {
                    "kind": "literal",
                      "value": [
                          {{.Key}},
                      ]
                  }
              {{- end }}
      }
  ]
}
`
)