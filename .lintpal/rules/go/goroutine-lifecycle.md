---
severity: high
threshold: 0.96
title: Bound goroutine lifetimes
---

Changed Go code that starts a goroutine must provide a clear termination path and prevent the goroutine from leaking after its owning operation ends.
