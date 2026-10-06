package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"icmongolang/internal/modules/users"
	usermodels "icmongolang/internal/modules/users/models"
	"icmongolang/internal/modules/users/presenter"
	"icmongolang/pkg/httpErrors"
	"icmongolang/pkg/responses"
	"icmongolang/pkg/utils"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// CreateFile – POST /user/files
// @Security     BearerAuth
// @Summary      Create user file
// @Tags         user-files
// @Accept       json
// @Produce      json
// @Param        body body presenter.FileCreate true "File data"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=usermodels.SdUserFile}
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/files [post]
func (h *userHandler) CreateFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(presenter.FileCreate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		if err := utils.ValidateStruct(r.Context(), req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		f := &usermodels.SdUserFile{
			FileName:   req.FileName,
			FileType:   req.FileType,
			FilePath:   req.FilePath,
			FileTypeID: req.FileTypeID,
			UID:        &req.UID,
			Status:     req.Status,
		}
		out, err := h.fileUC.CreateFile(r.Context(), f)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(out))
	}
}

// GetFile – GET /user/files/{id}
// @Security     BearerAuth
// @Summary      Get file by ID
// @Tags         user-files
// @Produce      json
// @Param        id path int true "File ID"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=usermodels.SdUserFile}
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Router       /user/files/{id} [get]
func (h *userHandler) GetFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconvAtoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		out, err := h.fileUC.GetFile(r.Context(), id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(out))
	}
}

// UpdateFile – PUT /user/files/{id}
// @Security     BearerAuth
// @Summary      Update file
// @Tags         user-files
// @Accept       json
// @Produce      json
// @Param        id path int true "File ID"
// @Param        body body presenter.FileUpdate true "Fields to update"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=usermodels.SdUserFile}
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/files/{id} [put]
func (h *userHandler) UpdateFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconvAtoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		req := new(presenter.FileUpdate)
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		values := map[string]interface{}{}
		if req.FileName != nil { values["file_name"] = *req.FileName }
		if req.FilePath != nil { values["file_path"] = *req.FilePath }
		if req.Status != nil   { values["status"] = *req.Status }
		out, err := h.fileUC.UpdateFile(r.Context(), id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(out))
	}
}

// DeleteFile – DELETE /user/files/{id}
// @Security     BearerAuth
// @Summary      Delete file
// @Tags         user-files
// @Param        id path int true "File ID"
// @Success      200  {object}  responses.SwaggerSuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Router       /user/files/{id} [delete]
func (h *userHandler) DeleteFile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconvAtoi(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err))); return
		}
		if err := h.fileUC.DeleteFile(r.Context(), id); err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(struct{}{}))
	}
}

// ListFiles – GET /user/files
// @Security     BearerAuth
// @Summary      List user files
// @Tags         user-files
// @Produce      json
// @Param        page         query int    false "Page" default(1)
// @Param        per_page     query int    false "Items per page" default(20)
// @Param        uid          query string false "Filter by uid"
// @Param        file_type_id query int    false "Filter by file_type_id"
// @Param        status       query string false "Status csv"
// @Success      200  {object}  responses.SwaggerSuccessResponse{data=presenter.FileListResponse}
// @Router       /user/files [get]
func (h *userHandler) ListFiles() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		page, perPage := parsePageParams(q)
		f := users.FileFilter{Limit: perPage, Offset: (page - 1) * perPage}
		if v := q.Get("uid"); v != "" { f.UID = &v }
		if v := q.Get("file_type_id"); v != "" {
			if n, err := strconvAtoi(v); err == nil { f.FileTypeID = &n }
		}
		if v := q.Get("status"); v != "" {
			for _, s := range strings.Split(v, ",") {
				if n, err := strconvAtoi16(strings.TrimSpace(s)); err == nil {
					f.Statuses = append(f.Statuses, n)
				}
			}
		}
		items, total, err := h.fileUC.ListFiles(r.Context(), f)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err)); return
		}
		list := make([]*presenter.FileResponse, 0, len(items))
		for _, it := range items {
			list = append(list, &presenter.FileResponse{
				ID:         it.ID,
				FileName:   it.FileName,
				FileType:   it.FileType,
				FilePath:   it.FilePath,
				FileTypeID: it.FileTypeID,
				UID:        derefStr(it.UID),
				FileDate:   it.FileDate,
				Status:     it.Status,
			})
		}
		render.Respond(w, r, responses.CreateSuccessResponse(
			presenter.NewFileListResponse(list, page, perPage, total),
		))
	}
}