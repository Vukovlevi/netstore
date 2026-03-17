<?php
function getDbConnection() {
    $env = parse_ini_file(__DIR__ . "/.env");
    $db = new mysqli($env['DB_HOST'], $env['DB_USER'], $env['DB_PASSWORD'], $env['DB_NAME']);
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