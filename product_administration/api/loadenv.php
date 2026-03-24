<?php
function loadEnv() {
    if (!file_exists('./.env')) {
        return;
    }
    $env = parse_ini_file('./.env');
    foreach ($env as $key => $value) {
        if (getenv($key) === false) {
            putenv("$key=$value");
        }
    }
}