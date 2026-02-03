---
title: "{{ replace .File.ContentBaseName "-" "." }}"
date: {{ .Date }}
version: "{{ replace .File.ContentBaseName "v" "" | replace "-" "." }}"
---

## What's New

-
