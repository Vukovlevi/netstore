<?php

use PHPUnit\Framework\TestCase;

abstract class HandlerTestCase extends TestCase
{
    protected function setUp(): void
    {
        DbStub::reset();
        $_GET = [];
        $_POST = [];
        $_REQUEST = [];
        $_COOKIE = [];
        http_response_code(200);
    }

    protected function runAndCapture(callable $fn): array
    {
        ob_start();
        try {
            $fn();
        } finally {
            $body = ob_get_clean();
        }
        return [
            'body' => $body,
            'json' => json_decode($body, true),
            'code' => http_response_code(),
        ];
    }
}
