"use client";

import axios from "axios";
import { schemas } from "package-ts";
import { useEffect, useState } from "react";
import { z } from "zod";

type System = z.infer<typeof schemas.model_System>;

const API_BASE_URL = "http://localhost:3003";

// システム一覧のレスポンススキーマ
const SystemsResponseSchema = z.array(schemas.model_System);

export default function Home() {
  const [systems, setSystems] = useState<System[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchSystems = async () => {
      try {
        setLoading(true);
        const response = await axios.get(`${API_BASE_URL}/api/v1/systems`);

        // zodスキーマでレスポンスをバリデーション
        const validatedData = SystemsResponseSchema.parse(response.data);
        setSystems(validatedData);
        setError(null);
      } catch (err) {
        console.error("システム一覧の取得に失敗しました:", err);

        // zodバリデーションエラーかAPIエラーかを判別
        if (err instanceof z.ZodError) {
          setError("APIレスポンスの形式が正しくありません");
          console.error("バリデーションエラー:", err.errors);
        } else {
          setError("システム一覧の取得に失敗しました");
        }
      } finally {
        setLoading(false);
      }
    };

    fetchSystems();
  }, []);

  if (loading) {
    return <div>読み込み中...</div>;
  }

  if (error) {
    return <div style={{ color: "red" }}>{error}</div>;
  }

  return (
    <div>
      <h1>システム一覧</h1>

      {systems.length === 0 ? (
        <p>システムが登録されていません</p>
      ) : (
        <table border={1} style={{ borderCollapse: "collapse", width: "100%" }}>
          <thead>
            <tr>
              <th>システム名</th>
              <th>自治体ID</th>
              <th>メールアドレス</th>
              <th>電話番号</th>
              <th>作成日時</th>
              <th>更新日時</th>
              <th>備考</th>
            </tr>
          </thead>
          <tbody>
            {systems.map((system) => (
              <tr key={system.id}>
                <td>{system.systemName}</td>
                <td>{system.localGovernmentId || "-"}</td>
                <td>{system.mailAddress}</td>
                <td>{system.telephone || "-"}</td>
                <td>{new Date(system.createdAt).toLocaleString("ja-JP")}</td>
                <td>{new Date(system.updatedAt).toLocaleString("ja-JP")}</td>
                <td>{system.remark || "-"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {systems.length > 0 && (
        <p>{systems.length}件のシステムが登録されています</p>
      )}
    </div>
  );
}
