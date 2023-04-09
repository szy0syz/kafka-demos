# mview-gin

## Notes

```bash
curl -X POST -H "Content-Type: application/json" -d '{"candidate": "Jerry Shi"}' http://localhost:3000/vote
```

<img width="1487" alt="image" src="https://user-images.githubusercontent.com/10555820/230775405-cf497725-214d-44fb-877c-d3799eec5f8d.png">

> 在这个示例中，Materialized View 的概念体现在我们对投票数据的实时统计和存储。具体来说，common.VoteCounts 变量扮演了 Materialized View 的角色。每当有新的投票事件进入 Kafka 消费者时，我们会对该候选人的票数进行更新。这样，我们始终维护着一个实时更新的票数统计。
>
> 要进一步强调 Materialized View 的概念，我们可以定期将 common.VoteCounts 数据持久化到磁盘或数据库。这样，即使在应用程序重新启动后，我们也可以从上次离线的地方开始，重新构建实时的票数统计。
