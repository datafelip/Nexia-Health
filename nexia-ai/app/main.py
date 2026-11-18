
from datetime import date, datetime, timezone
from pathlib import Path
from typing import Literal

import joblib
import pandas as pd

from fastapi import FastAPI, HTTPException
from pydantic import BaseModel, Field


app = FastAPI(
    title="Nexia AI",
    description=(
        "Serviço de Inteligência Artificial do Nexia Health. "
        "Modelo experimental treinado com dados sintéticos."
    ),
    version="1.1.0",
)


# Localização dos arquivos do projeto
BASE_DIR = Path(__file__).resolve().parents[2]

MODEL_PATH = (
    BASE_DIR
    / "nexia-ml"
    / "models"
    / "modelo_demanda.joblib"
)


class PrevisaoRequest(BaseModel):
    data: date = Field(
        description="Data para a previsão no formato AAAA-MM-DD"
    )

    feriado: Literal[0, 1] = Field(
        default=0,
        description=(
            "Indicador de feriado: 0 = não, 1 = sim. "
            "No dataset atual, esse indicador é sintético."
        ),
    )


class PrevisaoResponse(BaseModel):
    data: date
    atendimentos_previstos: int
    modelo: str
    origem_dados: str
    observacao: str


def carregar_modelo():
    if not MODEL_PATH.is_file():
        raise HTTPException(
            status_code=503,
            detail="Modelo de Machine Learning não encontrado.",
        )

    try:
        artefato = joblib.load(MODEL_PATH)
    except Exception as erro:
        raise HTTPException(
            status_code=503,
            detail="Não foi possível carregar o modelo.",
        ) from erro

    return artefato


@app.get("/health", tags=["Sistema"])
def health_check():
    return {
        "status": "online",
        "service": "nexia-ai",
        "version": "1.1.0",
        "timestamp": datetime.now(timezone.utc).isoformat(),
    }


@app.get("/api/v1/model/status", tags=["Machine Learning"])
def model_status():
    disponivel = MODEL_PATH.is_file()

    return {
        "model": "hospital-demand-forecast",
        "status": "trained" if disponivel else "not_trained",
        "trained": disponivel,
        "dataset": "synthetic",
        "message": (
            "Artefato do modelo disponível."
            if disponivel
            else "O modelo ainda não foi treinado."
        ),
    }


@app.post(
    "/api/v1/predict",
    response_model=PrevisaoResponse,
    tags=["Machine Learning"],
)
def prever_demanda(dados: PrevisaoRequest):
    artefato = carregar_modelo()

    modelo = artefato["modelo"]
    features = artefato["features"]

    dia_semana = dados.data.weekday()

    entrada = pd.DataFrame(
        [
            {
                "dia_semana": dia_semana,
                "mes": dados.data.month,
                "fim_semana": int(dia_semana >= 5),
                "feriado": dados.feriado,
            }
        ]
    )

    entrada = entrada[features]

    previsao = modelo.predict(entrada)[0]

    quantidade = max(0, round(float(previsao)))

    return PrevisaoResponse(
        data=dados.data,
        atendimentos_previstos=quantidade,
        modelo="RandomForestRegressor",
        origem_dados="sinteticos",
        observacao=(
            "Previsão experimental baseada em dados sintéticos. "
            "Não utilizar para decisões clínicas ou dimensionamento "
            "real de equipes hospitalares."
        ),
    )
