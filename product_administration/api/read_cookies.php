<?php
function authentication() {
  if(!isset($_COOKIE['auth_token'])) {
      return false;
    } else {
      $cookieCheck = getData("SELECT user.id, user.username, role.name FROM session INNER JOIN user ON session.user_id=user.id INNER JOIN role ON user.role_id=role.id WHERE token = ? AND expires_at > NOW()", "s", array($_COOKIE['auth_token']));
      if(empty($cookieCheck)) {
        return false;
      }

      $_REQUEST['user'] = $cookieCheck[0];
    }
    return true;
}