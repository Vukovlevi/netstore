<?php
require './sql_functions.php';
require './read_cookies.php';
$method = $_SERVER['REQUEST_METHOD'];
$uri = explode('/', parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH));
$body = json_decode(file_get_contents('php://input'), true);

authentication();
switch(end($uri)) {
    case "category":
        if ($method == "GET") {
            if (isset($_GET["id"])) {
                $id = $_GET["id"];
                $sql = "SELECT id, name FROM category WHERE id = ? AND deleted_at IS NULL";
                $category = getData($sql, "i", [$id]);
                if (empty($category)) {
                    echo json_encode(["message" => "Kategória nem található vagy törölve lett!"], JSON_UNESCAPED_UNICODE);
                    return http_response_code(404);
                }
                echo json_encode($category[0], JSON_UNESCAPED_UNICODE);
            } else {
                $sql = "SELECT id, name FROM category WHERE deleted_at IS NULL ORDER BY name";
                $categories = getData($sql);
                
                echo json_encode($categories, JSON_UNESCAPED_UNICODE);
            }
        }
        else if ($method == "POST") {
            if (empty($body["name"])) {
                echo json_encode(["message" => "Hiányzó adat!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }
            $name = $body["name"];
            $checkSql = "SELECT id, deleted_at FROM category WHERE name = ?";
            $existing = getData($checkSql, "s", [$name]);
            if (!empty($existing)) {
                if (is_null($existing[0]["deleted_at"])) {
                    echo json_encode(["message" => "Már létezik ilyen nevű aktív kategória!"], JSON_UNESCAPED_UNICODE);
                    return http_response_code(409);
                } else {
                    echo json_encode(["message" => "Már létezik ilyen nevű kategória, de törölve lett. Kérem, használjon más nevet!"], JSON_UNESCAPED_UNICODE);
                    return http_response_code(409);
                }
            }
            $insertSql = "INSERT INTO category (name) VALUES (?)";
            $success = changeData($insertSql, "s", [$name]);
            if (!$success) {
                echo json_encode(["message" => "Sikertelen művelet, adatbázis hiba!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(500);
            }
            $lastIdSql = "SELECT id, name FROM category WHERE name = ? AND deleted_at IS NULL";
            $newCategory = getData($lastIdSql, "s", [$name]);
            
            echo json_encode($newCategory[0], JSON_UNESCAPED_UNICODE);
            http_response_code(201);
        }
        else if ($method == "PUT") {
            if (empty($body["id"]) || empty($body["name"])) {
                echo json_encode(["message" => "Hiányzó adat: 'id' és 'name' megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }
            $id = $body["id"];
            $name = $body["name"];
            $checkSql = "SELECT id, name FROM category WHERE id = ? AND deleted_at IS NULL";
            $categoryToUpdate = getData($checkSql, "i", [$id]);
            if (empty($categoryToUpdate)) {
                echo json_encode(["message" => "Nincs ilyen azonosítójú kategória, vagy törölve lett!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(404);
            }
            $checkNameSql = "SELECT id FROM category WHERE name = ? AND id != ?";
            $existingName = getData($checkNameSql, "si", [$name, $id]);
            if (!empty($existingName)) {
                echo json_encode(["message" => "Ez a név már foglalt egy másik (esetleg törölt) kategória által!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(409);
            }
            $updateSql = "UPDATE category SET name = ? WHERE id = ?";
            $success = changeData($updateSql, "si", [$name, $id]);
            if (!$success) {
                echo json_encode(["message" => "Sikertelen művelet, adatbázis hiba!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(500);
            }
            $updatedCategory = getData($checkSql, "i", [$id]);
            echo json_encode($updatedCategory[0], JSON_UNESCAPED_UNICODE);
        }
        else if ($method == "DELETE") {
            if (empty($_GET["id"])) {
                echo json_encode(["message" => "Hiányzó adat: 'id' query paraméter megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }
            $id = $_GET["id"];
            $checkSql = "SELECT id FROM category WHERE id = ? AND deleted_at IS NULL";
            $categoryToDelete = getData($checkSql, "i", [$id]);
            if (empty($categoryToDelete)) {
                echo json_encode(["message" => "Nincs ilyen azonosítójú kategória, vagy már törölve lett!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(404);
            }
            $checkSubCatSql = "SELECT id FROM sub_category WHERE category_id = ? AND deleted_at IS NULL";
            $subCategories = getData($checkSubCatSql, "i", [$id]);
            
            if (!empty($subCategories)) {
                echo json_encode([
                    "message" => "Törlés sikertelen: Ez a kategória aktív alkategóriá(ka)t tartalmaz. Előbb törölje az alkategóriákat."
                ], JSON_UNESCAPED_UNICODE);
                return http_response_code(409);
            }
            $deleteSql = "UPDATE category SET deleted_at = CURDATE() WHERE id = ?";
            $success = changeData($deleteSql, "i", [$id]);
            if (!$success) {
                echo json_encode(["message" => "Sikertelen művelet, adatbázis hiba!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(500);
            }
            
            http_response_code(204);
        }
        else {
            return http_response_code(405);
        }
        break;
    default:
        return http_response_code(404);
}