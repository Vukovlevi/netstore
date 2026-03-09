<?php
function getData($operation, $types = null, $data = null) {
    $db = new mysqli('localhost', 'root', 'root', 'netstore');
    if ($db->connect_errno != 0) {
        return $db->connect_error;
    }
    
    $db->set_charset("utf8mb4");

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
    $db = new mysqli('localhost', 'root', 'root', 'netstore');
    if ($db->connect_errno != 0) {
        return $db->connect_error;
    }
    
    $db->set_charset("utf8mb4");

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