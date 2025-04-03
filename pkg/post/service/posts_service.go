package service

import (
	template "github.com/axlle-com/blog/pkg/template/provider"
	user "github.com/axlle-com/blog/pkg/user/provider"
	"sync"

	"github.com/axlle-com/blog/pkg/alias"
	"github.com/axlle-com/blog/pkg/app/logger"
	"github.com/axlle-com/blog/pkg/app/models/contracts"
	"github.com/axlle-com/blog/pkg/file/provider"
	gallery "github.com/axlle-com/blog/pkg/gallery/provider"
	"github.com/axlle-com/blog/pkg/post/models"
	"github.com/axlle-com/blog/pkg/post/repository"
)

type PostsService struct {
	postRepo          repository.PostRepository
	categoriesService *CategoriesService
	categoryService   *CategoryService
	galleryProvider   gallery.GalleryProvider
	fileProvider      provider.FileProvider
	aliasProvider     alias.AliasProvider
	userProvider      user.UserProvider
	templateProvider  template.TemplateProvider
}

func NewPostsService(
	postRepo repository.PostRepository,
	categoriesService *CategoriesService,
	categoryService *CategoryService,
	galleryProvider gallery.GalleryProvider,
	fileProvider provider.FileProvider,
	aliasProvider alias.AliasProvider,
	userProvider user.UserProvider,
	templateProvider template.TemplateProvider,
) *PostsService {
	return &PostsService{
		postRepo:          postRepo,
		categoriesService: categoriesService,
		categoryService:   categoryService,
		galleryProvider:   galleryProvider,
		fileProvider:      fileProvider,
		aliasProvider:     aliasProvider,
		userProvider:      userProvider,
		templateProvider:  templateProvider,
	}
}

func (s *PostsService) GetAggregates(posts []*models.Post) []*models.Post {
	var templateIDs []uint
	var userIDs []uint
	var categoryIDs []uint

	templateIDsMap := make(map[uint]bool)
	userIDsMap := make(map[uint]bool)
	categoryIDsMap := make(map[uint]bool)

	for _, post := range posts {
		if post.TemplateID != nil {
			id := *post.TemplateID
			if !templateIDsMap[id] {
				templateIDs = append(templateIDs, id)
				templateIDsMap[id] = true
			}
		}
		if post.UserID != nil {
			id := *post.UserID
			if !userIDsMap[id] {
				userIDs = append(userIDs, id)
				userIDsMap[id] = true
			}
		}
		if post.PostCategoryID != nil {
			id := *post.PostCategoryID
			if !categoryIDsMap[id] {
				categoryIDs = append(categoryIDs, id)
				categoryIDsMap[id] = true
			}
		}
	}

	var wg sync.WaitGroup

	var users map[uint]contracts.User
	var templates map[uint]contracts.Template
	var categories map[uint]*models.PostCategory

	wg.Add(3)

	go func() {
		defer wg.Done()
		if len(templateIDs) > 0 {
			var err error
			templates, err = s.templateProvider.GetMapByIDs(templateIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		if len(userIDs) > 0 {
			var err error
			users, err = s.userProvider.GetMapByIDs(userIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	go func() {
		defer wg.Done()
		if len(categoryIDs) > 0 {
			var err error
			categories, err = s.categoriesService.GetMapByIDs(categoryIDs)
			if err != nil {
				logger.Error(err)
			}
		}
	}()

	wg.Wait()

	for _, post := range posts {
		if post.TemplateID != nil {
			post.Template = templates[*post.TemplateID]
		}
		if post.UserID != nil {
			post.User = users[*post.UserID]
		}
		if post.PostCategoryID != nil {
			post.Category = categories[*post.PostCategoryID]
		}
	}

	return posts
}

func (s *PostsService) WithPaginate(p contracts.Paginator, filter *models.PostFilter) ([]*models.Post, error) {
	return s.postRepo.WithPaginate(p, filter)
}
