<?php
function getData($operation, $types = null, $data = null) {
    $db = new mysqli('localhost', 'root', 'root', 'netstore');
    if ($db->connect_errno != 0) {
        return $db->connect_error;
    }

    if (!is_null($types) && !is_null($data)) {
        $stmt = $db->prepare($operation);
        $stmt->bind_param($types, ...$data);
        $stmt->execute();
        $eredmeny = $stmt->get_result();
    }
    else {
        $eredmeny = $db->query($operation);
    }
    
    if ($db->errno != 0) {
        return $db->error;
    }
    if ($eredmeny->num_rows == 0) {
        return [];
    }

    return $eredmeny->fetch_all(MYSQLI_ASSOC);
}

function changeData($operation, $types = null, $data = null) {
    $db = new mysqli('localhost', 'root', 'root', 'netstore');
    if ($db->connect_errno != 0) {
        return $db->connect_error;
    }

    if (!is_null($types) && !is_null($data)) {
        $stmt = $db->prepare($operation);
        $stmt->bind_param($types, ...$data);
        $stmt->execute();
    }
    else {
        $db->query($operation);
    }
    
    if ($db->errno != 0) {
        return $db->error;
    }

    return $db->affected_rows > 0 ? true : false;
}