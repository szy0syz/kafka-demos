package com.learnkafkastreams.topology;

import org.apache.kafka.common.serialization.Serdes;
import org.apache.kafka.streams.StreamsBuilder;
import org.apache.kafka.streams.Topology;
import org.apache.kafka.streams.kstream.Consumed;
import org.apache.kafka.streams.kstream.Printed;
import org.apache.kafka.streams.kstream.Produced;

public class GreetingsTopology {
  public static String GREETINGS = "greetings";

  public static String GREETINGS_UPPERCASE = "greetings_uppercase";
  public static String GREETINGS_CHINESE = "greetings_chinese";

  public static Topology buildTopology() {
    StreamsBuilder streamsBuilder = new StreamsBuilder();

    var greetingsStream = streamsBuilder
        .stream(GREETINGS,
            Consumed.with(Serdes.String(), Serdes.String()));

    var greetingsChineseStream = streamsBuilder
        .stream(GREETINGS_CHINESE,Consumed.with(Serdes.String(), Serdes.String()));

    var mergedStream = greetingsStream.merge(greetingsChineseStream);

    mergedStream
        .print(Printed.<String, String>toSysOut().withLabel("MergedStream"));

    var modifiedStream = mergedStream
        .mapValues(((readOnlyKey, value) -> value.toUpperCase()));

    modifiedStream
        .print(Printed.<String, String>toSysOut().withLabel("modifiedStream"));

    modifiedStream
        .to(GREETINGS_UPPERCASE
//            ,Produced.with(Serdes.String(), Serdes.String())
        );

    return streamsBuilder.build();
  }
}
