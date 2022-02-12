<h1 align="center">
  🏗 Render Go Template 🚧 
</h1>

<h3 align="center">
  Harness the power of <a href="https://pkg.go.dev/text/template">Go templates</a> in a lightweight CLI or GitHub Action

  [![GitHub Workflow Status](https://img.shields.io/github/workflow/status/YuKitsune/template-cli/CI)](https://github.com/YuKitsune/template-cli/actions?query=workflow:CI)
  [![Go Report Card](https://goreportcard.com/badge/github.com/YuKitsune/template-cli)](https://goreportcard.com/report/github.com/YuKitsune/template-cli)
  [![License](https://img.shields.io/github/license/YuKitsune/template-cli)](https://github.com/YuKitsune/template-cli/blob/main/LICENSE)
  [![Latest Release](https://img.shields.io/github/v/release/YuKitsune/template-cli?include_prereleases)](https://github.com/YuKitsune/template-cli/releases)

</h3>

# Quick start

## GitHub Action
```yaml
steps:
  - name: Render template files
    uses: yukitsune/template-cli@main
    with:
      args: --input ./templates/file1 --input ./templates/file2 \
        --value "person.name=Jason"\
        --value "person.age=${{ secrets.PERSON_AGE }}"\
        --value "secret=${{ secrets.GITHUB_TOKEN }}"
        --output .
```

## CLI
```shell
template -i ./templates/file1 -i ./templates/file2 \
  -v "person.name=Jason" \
  -v "person.age=${{ secrets.PERSON_AGE }}" \
  -v "secret=${{ secrets.GITHUB_TOKEN }}" \
  -o .
```

# Contributing

Contributions are what make the open source community such an amazing place to be, learn, inspire, and create.
Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`feature/AmazingFeature`)
3. Commit your Changes
4. Push to the Branch
5. Open a Pull Request
