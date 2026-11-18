
package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	ErrAIIndisponivel = errors.New("serviço de IA indisponível")
	ErrAIResposta     = errors.New("resposta inválida do serviço de IA")
)

type AIService struct {
	baseURL string
	client  *http.Client
}

type AIPrevisaoRequest struct {
	Data    string `json:"data"`
	Feriado int    `json:"feriado"`
}

type AIPrevisaoResponse struct {
	Data                 string `json:"data"`
	AtendimentosPrevistos int    `json:"atendimentos_previstos"`
	Modelo               string `json:"modelo"`
	OrigemDados          string `json:"origem_dados"`
	Observacao           string `json:"observacao"`
}

func NewAIService(baseURL string) *AIService {
	if strings.TrimSpace(baseURL) == "" {
		baseURL = "http://127.0.0.1:8000"
	}

	return &AIService{
		baseURL: strings.TrimRight(baseURL, "/"),
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *AIService) PreverDemanda(
	ctx context.Context,
	dados AIPrevisaoRequest,
) (*AIPrevisaoResponse, error) {

	// Validar a data antes de consultar o Python.
	data, err := time.Parse("2006-01-02", dados.Data)
	if err != nil || data.Format("2006-01-02") != dados.Data {
		return nil, fmt.Errorf(
			"data inválida: utilize o formato AAAA-MM-DD",
		)
	}

	// O modelo aceita apenas 0 ou 1 para feriado.
	if dados.Feriado != 0 && dados.Feriado != 1 {
		return nil, fmt.Errorf(
			"feriado deve ser 0 ou 1",
		)
	}

	corpo, err := json.Marshal(dados)
	if err != nil {
		return nil, fmt.Errorf(
			"erro ao preparar requisição: %w",
			err,
		)
	}

	url := s.baseURL + "/api/v1/predict"

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewReader(corpo),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"erro ao criar requisição de IA: %w",
			err,
		)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Enviar a requisição ao FastAPI.
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf(
			"%w: %v",
			ErrAIIndisponivel,
			err,
		)
	}
	defer resp.Body.Close()

	// Não permitir leitura ilimitada da resposta.
	leitor := io.LimitReader(resp.Body, 1024*1024)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"%w: HTTP %d",
			ErrAIIndisponivel,
			resp.StatusCode,
		)
	}

	var previsao AIPrevisaoResponse

	if err := json.NewDecoder(leitor).Decode(&previsao); err != nil {
		return nil, fmt.Errorf(
			"%w: %v",
			ErrAIResposta,
			err,
		)
	}

	// Validar a estrutura mínima da resposta.
	if previsao.Data != dados.Data ||
		previsao.Modelo == "" ||
		previsao.OrigemDados == "" ||
		previsao.AtendimentosPrevistos < 0 {

		return nil, ErrAIResposta
	}

	return &previsao, nil
}
