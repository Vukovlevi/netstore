<?php

final class DbStub
{
    public static array $getDataResponses = [];

    public static array $getDataCalls = [];

    public static array $changeDataResponses = [];

    public static array $changeDataCalls = [];

    public static function reset(): void
    {
        self::$getDataResponses = [];
        self::$getDataCalls = [];
        self::$changeDataResponses = [];
        self::$changeDataCalls = [];
    }

    public static function queueGetData(mixed ...$responses): void
    {
        foreach ($responses as $r) {
            self::$getDataResponses[] = $r;
        }
    }

    public static function queueChangeData(mixed ...$responses): void
    {
        foreach ($responses as $r) {
            self::$changeDataResponses[] = $r;
        }
    }
}
