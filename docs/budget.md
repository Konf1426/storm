# Budget Estimate (target <= 700 EUR)

This estimate uses public on-demand prices. Some prices are from third-party indexes that mirror the AWS Pricing API; verify with the AWS Pricing Calculator before production.

## Assumptions
- Provider: AWS
- Region: EU (Ireland) for EC2; US-East-1 for RDS/ElastiCache pricing reference
- Single-AZ for RDS and single-node ElastiCache
- ASG desired capacity: 1 (scale out increases cost linearly)
- 730 hours per month

## Cost Breakdown (approx, monthly)
- EC2 (t3a.medium, 1x): $0.0408/hr => ~$29.78/mo
- ALB: $0.0225/hr + $0.008 per LCU-hr (assume 1 LCU) => ~$22.27/mo
- RDS Postgres (db.t4g.micro): $0.03/hr => ~$21.90/mo
- ElastiCache Redis (cache.t4g.micro): $0.02/hr => ~$11.68/mo
- EBS (gp3) storage: $0.08 per GiB-month; 20 GB => ~$1.60/mo
- Data transfer, backups, logs: not included (varies by traffic)

Approx subtotal (excluding data transfer): ~$86-90/mo

## Notes
- ALB cost depends on LCU usage; heavy traffic increases LCU charges.
- RDS/ElastiCache prices vary by region and instance size; adjust in the AWS Pricing Calculator.
- To reach 100k connections / 500k msg/s, compute costs will increase; scale the EC2 linearly with ASG size.

## Scenarios (pre-remplis)
1) Minimal (dev/demo): 1x EC2, 1x RDS micro, 1x Redis micro, ALB minimal.
2) Medium: 2-3x EC2, RDS small, Redis small, ALB avec 2-3 LCU.
3) Haute charge: autoscaling >5 instances, DB/cache plus grands, couts > 700 EUR sans optimisations.

## Tableau de comparaison (a valider)
| Scenario | EC2 | DB | Redis | ALB/LCU | Estimation |
| --- | --- | --- | --- | --- | --- |
| Minimal | 1 x t3a.medium | db.t4g.micro | cache.t4g.micro | 1 LCU | ~90 USD/mo |
| Medium | 3 x t3a.medium | db.t4g.small | cache.t4g.small | 2 LCU | ~250-300 USD/mo |
| Haute charge | 6+ x t3a.medium | db.m6g.large | cache.m6g.large | 3+ LCU | > 700 USD/mo |
