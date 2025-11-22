<?php
require './sql_functions.php';
function authentication() {
    if(!isset($_COOKIE['auth_token'])) {
      echo "Cookie named '" . 'auth_token' . "' is not set!";
    } else {
      $cookieCheck = getData("SELECT id, user_id, token, expires_at FROM session WHERE token = ? AND expires_at > NOW()", "s", array($_COOKIE["auth_token"]));
    }
}