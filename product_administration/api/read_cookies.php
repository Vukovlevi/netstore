<?php
function authentication() {
    if(!isset($_COOKIE['auth_token'])) {
      //header("Location: http://localhost:8000/login");
      echo "asd";
      exit;
    } else {
      $cookieCheck = getData("SELECT id, user_id, token, expires_at FROM session WHERE token = ? AND expires_at > NOW()", "s", array($_COOKIE["auth_token"]));
      var_dump($cookieCheck);
      //if(empty($cookieCheck)) {
      //  header("Location: http://localhost:8000/login");
      //  exit;
      //}
    }
}