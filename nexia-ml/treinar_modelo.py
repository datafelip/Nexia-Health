
from pathlib import Path
import json

import joblib
import numpy as np
import pandas as pd

from sklearn.ensemble import RandomForestRegressor
from sklearn.metrics import (
    mean_absolute_error,
    mean_squared_error,
    r2_score,
)


BASE_DIR = Path(__file__).resolve().parent

DATASET_PATH = (
    BASE_DIR / "datasets" / "atendimentos_sinteticos.csv"
)

MODEL_PATH = (
    BASE_DIR / "models" / "modelo_demanda.joblib"
)

METRICS_PATH = (
    BASE_DIR / "models" / "metricas.json"
)

FEATURES = [
    "dia_semana",
    "mes",
    "fim_semana",
    "feriado",
]


def calcular_metricas(y_real, y_previsto):
    return {
        "mae": round(
            float(mean_absolute_error(y_real, y_previsto)), 4
        ),
        "rmse": round(
            float(np.sqrt(mean_squared_error(y_real, y_previsto))), 4
        ),
        "r2": round(
            float(r2_score(y_real, y_previsto)), 4
        ),
    }


def treinar_modelo():
    print("=" * 60)
    print("NEXIA AI - TREINAMENTO DE MACHINE LEARNING")
    print("=" * 60)

    if not DATASET_PATH.exists():
        raise FileNotFoundError(
            f"Dataset não encontrado: {DATASET_PATH}"
        )

    # 1. Carregar o histórico sintético
    df = pd.read_csv(
        DATASET_PATH,
        parse_dates=["data"]
    )

    df = df.sort_values("data").reset_index(drop=True)

    if df.empty or len(df) < 30:
        raise ValueError(
            "O dataset precisa conter pelo menos 30 registros."
        )

    if df[FEATURES + ["total_atendimentos"]].isna().any().any():
        raise ValueError(
            "O dataset contém valores ausentes."
        )

    # 2. Separar os dados cronologicamente
    ponto_corte = int(len(df) * 0.8)

    treino = df.iloc[:ponto_corte]
    teste = df.iloc[ponto_corte:]

    X_treino = treino[FEATURES]
    y_treino = treino["total_atendimentos"]

    X_teste = teste[FEATURES]
    y_teste = teste["total_atendimentos"]

    print(f"Registros totais: {len(df)}")
    print(f"Treinamento: {len(treino)} registros")
    print(f"Teste: {len(teste)} registros")

    print(
        "Período de treinamento:",
        treino["data"].min().date(),
        "até",
        treino["data"].max().date()
    )

    print(
        "Período de teste:",
        teste["data"].min().date(),
        "até",
        teste["data"].max().date()
    )

    # 3. Criar e treinar o modelo
    modelo = RandomForestRegressor(
        n_estimators=200,
        max_depth=10,
        min_samples_leaf=3,
        random_state=42,
        n_jobs=-1,
    )

    modelo.fit(X_treino, y_treino)

    # 4. Realizar previsões no conjunto de teste
    previsoes = modelo.predict(X_teste)

    metricas_modelo = calcular_metricas(
        y_teste,
        previsoes
    )

    # 5. Comparar com uma previsão de referência
    media_treinamento = float(y_treino.mean())

    previsoes_referencia = np.full(
        shape=len(y_teste),
        fill_value=media_treinamento,
    )

    metricas_referencia = calcular_metricas(
        y_teste,
        previsoes_referencia
    )

    # 6. Salvar o modelo e suas informações
    MODEL_PATH.parent.mkdir(
        parents=True,
        exist_ok=True
    )

    artefato = {
        "modelo": modelo,
        "features": FEATURES,
        "tipo": "RandomForestRegressor",
        "dataset": "sintetico",
        "periodo_treinamento": {
            "inicio": str(treino["data"].min().date()),
            "fim": str(treino["data"].max().date()),
        },
    }

    joblib.dump(artefato, MODEL_PATH)

    relatorio = {
        "nome": "Nexia AI - Previsão de Demanda",
        "modelo": "RandomForestRegressor",
        "origem_dados": "sinteticos",
        "total_registros": len(df),
        "registros_treinamento": len(treino),
        "registros_teste": len(teste),
        "features": FEATURES,
        "metricas_modelo": metricas_modelo,
        "metricas_referencia": metricas_referencia,
    }

    with open(
        METRICS_PATH,
        "w",
        encoding="utf-8"
    ) as arquivo:
        json.dump(
            relatorio,
            arquivo,
            indent=4,
            ensure_ascii=False
        )

    # 7. Exibir os resultados
    print("\nRESULTADOS DO MODELO")
    print("-" * 60)

    print(
        f"MAE: {metricas_modelo['mae']:.4f} atendimentos"
    )

    print(
        f"RMSE: {metricas_modelo['rmse']:.4f} atendimentos"
    )

    print(
        f"R²: {metricas_modelo['r2']:.4f}"
    )

    print("\nMODELO DE REFERÊNCIA (MÉDIA HISTÓRICA)")
    print("-" * 60)

    print(
        f"MAE: {metricas_referencia['mae']:.4f} atendimentos"
    )

    print(
        f"RMSE: {metricas_referencia['rmse']:.4f} atendimentos"
    )

    print(
        f"R²: {metricas_referencia['r2']:.4f}"
    )

    print("\nARQUIVOS GERADOS")
    print("-" * 60)
    print(f"Modelo: {MODEL_PATH}")
    print(f"Métricas: {METRICS_PATH}")

    print("\nTreinamento concluído com sucesso!")


if __name__ == "__main__":
    treinar_modelo()
