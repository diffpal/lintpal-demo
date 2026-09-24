---
severity: medium
threshold: 0.95
title: Propagate request context
---

Changed Go code must propagate the caller's context through request-scoped operations instead of replacing it with a background context.
