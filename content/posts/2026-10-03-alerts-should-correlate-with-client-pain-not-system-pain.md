---
title: 'Alerts should correlate with client pain, not system pain'

date: '2026-10-03 09:34:54.467009 +0000 UTC'
dateFormat: '02 Jan, 2006'

draft: true
---

We define alerts to notify us when there are issues affecting our clients. These are especially important when there are 
issues that cannot wait to be resolved.

However, alerting on the wrong thing can burn out on-call engineers, erode trust in our alerts, delay incident response,
and even cause us to miss genuine issues affecting clients.

> So, how can we define effective alerts?

We should define alerts that correlate with client pain, not system pain. 

## 1. Measure symptoms, not causes.

Our systems can fail due to a number of causes. But, not all of those causes result in client-facing impact. Alerting on
causes would ...

On the other hand, symptoms, by definition, are client-facing. Alerting on symptoms ensures that the (precious) time of 
an on-call engineer is spent on responding, diagnosing, and treating the most important issues.

## 2. Alert on breached SLOs.

what are SLOs? why does it make sense to alert on breaches?

## 3. State impact using a DSL.

what is a DSL? why is a DSL useful here?