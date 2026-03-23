<?php
function loadEnv() {
    $env = parse_ini_file('./.env');
    foreach ($env as $key => $value) {
        if (!isset($_ENV[$key])) {
            $_ENV[$key] = $value;
        }
    }
}