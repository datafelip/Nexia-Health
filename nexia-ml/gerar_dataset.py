
from pathlib import Path

import numpy as np
import pandas as pd


# Configurações do dataset
SEMENTE = 42
DATA_INICIAL = "2024-01-01"
DATA_FINAL = "2026-10-31"

DIRETORIO_ATUAL = Path(__file__).resolve().parent
CAMINHO_DATASET = (
    DIRETORIO_ATUAL / "datasets" / "atendimentos_sinteticos.csv"
)


def gerar_dataset():
    gerador = np.random.default_rng(SEMENTE)

    datas = pd.date_range(
        start=DATA_INICIAL,
        end=DATA_FINAL,
        freq="D"
    )

    registros = []

    for data in datas:
        dia_semana = data.weekday()
        mes = data.month

        # Indicador de fim de semana
        fim_semana = int(dia_semana >= 5)

        # Indicador fictício de feriado
        # Não corresponde ao calendário oficial brasileiro.
        feriado = int(gerador.random() < 0.03)

        # Demanda diária simulada
        demanda_base = 45

        # Variação conforme o dia da semana
        efeito_semana = {
            0: 12,
            1: 8,
            2: 5,
            3: 4,
            4: 7,
            5: -6,
            6: -9
        }

        # Sazonalidade hipotética
        efeito_sazonal = 8 if mes in [5, 6, 7, 8] else 0

        # Variação aleatória
        ruido = gerador.normal(loc=0, scale=5)

        total_atendimentos = (
            demanda_base
            + efeito_semana[dia_semana]
            + efeito_sazonal
            - (5 * feriado)
            + ruido
        )

        total_atendimentos = max(
            0,
            round(total_atendimentos)
        )

        registros.append({
            "data": data.strftime("%Y-%m-%d"),
            "dia_semana": dia_semana,
            "mes": mes,
            "fim_semana": fim_semana,
            "feriado": feriado,
            "total_atendimentos": total_atendimentos
        })

    dataframe = pd.DataFrame(registros)

    CAMINHO_DATASET.parent.mkdir(
        parents=True,
        exist_ok=True
    )

    dataframe.to_csv(
        CAMINHO_DATASET,
        index=False,
        encoding="utf-8"
    )

    print("=" * 55)
    print("NEXIA AI - GERADOR DE DATASET SINTÉTICO")
    print("=" * 55)
    print(f"Registros gerados: {len(dataframe)}")
    print(f"Período: {DATA_INICIAL} até {DATA_FINAL}")
    print(f"Média diária: {dataframe['total_atendimentos'].mean():.2f}")
    print(f"Mínimo: {dataframe['total_atendimentos'].min()}")
    print(f"Máximo: {dataframe['total_atendimentos'].max()}")
    print(f"Arquivo salvo em: {CAMINHO_DATASET}")
    print("=" * 55)

    print("\nPrimeiros cinco registros:")
    print(dataframe.head().to_string(index=False))


if __name__ == "__main__":
    gerar_dataset()
