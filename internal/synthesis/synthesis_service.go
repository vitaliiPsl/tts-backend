package synthesis

import (
	"encoding/json"
	"vitaliiPsl/synthesizer/internal/history"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/model"
	"vitaliiPsl/synthesizer/internal/requests"

	"github.com/gofiber/fiber/v2"
)

type SynthesisService interface {
	HandleSynthesisRequest(req *requests.SynthesisRequest, userId string) (*SynthesisResponse, error)
}

type SynthesisServiceImpl struct {
	modelService   model.ModelService
	historyService history.HistoryService
}

func NewSynthesisService(modelService model.ModelService, historyService history.HistoryService) *SynthesisServiceImpl {
	return &SynthesisServiceImpl{
		modelService:   modelService,
		historyService: historyService,
	}
}

func (s *SynthesisServiceImpl) HandleSynthesisRequest(req *requests.SynthesisRequest, userId string) (*SynthesisResponse, error) {
	logger.Logger.Info("Handling synthesis...", "userId", userId)

	model, err := s.modelService.GetModelById(req.ModelId)
	if err != nil {
		return nil, err
	}

	response, err := s.performSynthesis(model, req.Text)
	if err != nil {
		return nil, err
	}

	if userId != "" {
		err = s.saveHistoryRecord(req, userId)
		if err != nil {
			return nil, err
		}
	}

	logger.Logger.Info("Handled synthesis.", "userId", userId)
	return response, nil
}

func (s *SynthesisServiceImpl) performSynthesis(model *model.ModelDto, text string) (*SynthesisResponse, error) {
	logger.Logger.Info("Performing synthesis...", "name", model.Name, "language", model.Name, "url", model.Url)

	agent := fiber.Post(model.Url)
	agent.JSON(fiber.Map{"text": text})

	statusCode, resBody, errs := agent.Bytes()
	if len(errs) > 0 || statusCode != 200 {
		logger.Logger.Error("Failed to synthesize speech")
		return nil, errs[0]
	}

	var response *SynthesisResponse
	err := json.Unmarshal(resBody, &response)
	if err != nil {
		logger.Logger.Error("Failed to unmarshal synthesis response")
		return nil, err
	}

	return response, nil
}

func (s *SynthesisServiceImpl) saveHistoryRecord(req *requests.SynthesisRequest, userId string) error {
	historyDto := &history.HistoryRecordDto{
		UserId: userId,
		Text:   req.Text,
	}

	_, err := s.historyService.SaveHistoryRecord(historyDto)
	return err
}
