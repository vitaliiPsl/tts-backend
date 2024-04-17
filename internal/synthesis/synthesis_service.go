package synthesis

import (
	"encoding/json"
	"os"
	"vitaliiPsl/synthesizer/internal/history"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/requests"

	"github.com/gofiber/fiber/v2"
)

type SynthesisService struct {
	synthesisServiceUrl string
	historyService      *history.HistoryService
}

func NewSynthesisService(historyService *history.HistoryService) *SynthesisService {
	synthesisServiceUrl := os.Getenv("SYNTHESIS_SERVICE_URL")
	return &SynthesisService{
		synthesisServiceUrl: synthesisServiceUrl,
		historyService:      historyService,
	}
}

func (s *SynthesisService) HandleSynthesisRequest(req *requests.SynthesisRequest, userId string) (*SynthesisResponse, error) {
	logger.Logger.Info("Handling synthesis...", "userId", userId)

	response, err := s.performSynthesis(req)
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

func (s *SynthesisService) performSynthesis(req *requests.SynthesisRequest) (*SynthesisResponse, error) {
	agent := fiber.Post(s.synthesisServiceUrl)
	agent.JSON(req)
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

func (s *SynthesisService) saveHistoryRecord(req *requests.SynthesisRequest, userId string) error {
	historyDto := &history.HistoryRecordDto{
		UserId: userId,
		Text:   req.Text,
	}

	_, err := s.historyService.SaveHistoryRecord(historyDto)
	return err
}
