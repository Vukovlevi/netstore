<?php
require './loadenv.php';
function getDbConnection() {
    loadEnv();
    $db = new mysqli($_ENV['DB_HOST'], $_ENV['DB_USER'], $_ENV['DB_PASSWORD'], $_ENV['DB_NAME']);
    if ($db->connect_errno != 0) {
        return null;
    }
    $db->set_charset("utf8mb4");
    return $db;
}

function getData($operation, $types = null, $data = null) {
    $db = getDbConnection();
    if ($db === null) {
        return "Database connection failed";
    }

    if (!is_null($types) && !is_null($data)) {
        $stmt = $db->prepare($operation);
        $stmt->bind_param($types, ...$data);
        $stmt->execute();
        $result = $stmt->get_result();
    }
    else {
        $result = $db->query($operation);
    }

    if ($db->errno != 0) {
        return $db->error;
    }

    if ($result === true || $result === false) {
        return [];
    }
    if ($result->num_rows == 0) {
        return [];
    }

    return $result->fetch_all(MYSQLI_ASSOC);
}

function changeData($operation, $types = null, $data = null) {
    $db = getDbConnection();
    if ($db === null) {
        return false;
    }

    if (!is_null($types) && !is_null($data)) {
        $stmt = $db->prepare($operation);
        $stmt->bind_param($types, ...$data);
        $success = $stmt->execute();
        $stmt->close();
    }
    else {
        $success = $db->query($operation);
    }

    if ($db->errno != 0) {
        return false;
    }

    return $success;
}
?>