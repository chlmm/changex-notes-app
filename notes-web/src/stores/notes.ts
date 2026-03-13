import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { BookNote, VideoNote, KnowledgeNote, SkillNote, ProblemNote, IndexNote, QuoteNote, GitHubRepoNote, AnimeNote, MovieNote, GameNote, NoteStats } from '@/types/note'

// API 模块
import * as booksApi from '@/api/books'
import * as videosApi from '@/api/videos'
import * as knowledgeApi from '@/api/knowledge'
import * as skillsApi from '@/api/skills'
import * as problemsApi from '@/api/problems'
import * as indexApi from '@/api/indexNotes'
import * as quotesApi from '@/api/quotes'
import * as githubApi from '@/api/github'
import * as animeApi from '@/api/anime'
import * as moviesApi from '@/api/movies'
import * as gamesApi from '@/api/games'
import * as searchApi from '@/api/search'

export const useNotesStore = defineStore('notes', () => {
  const books = ref<BookNote[]>([])
  const videos = ref<VideoNote[]>([])
  const knowledge = ref<KnowledgeNote[]>([])
  const skills = ref<SkillNote[]>([])
  const problems = ref<ProblemNote[]>([])
  const indexNotes = ref<IndexNote[]>([])
  const quotes = ref<QuoteNote[]>([])
  const githubRepos = ref<GitHubRepoNote[]>([])
  const anime = ref<AnimeNote[]>([])
  const movies = ref<MovieNote[]>([])
  const games = ref<GameNote[]>([])
  const loading = ref(false)
  const loaded = ref(false)

  const stats = computed<NoteStats>(() => ({
    totalBooks: books.value.length,
    totalVideos: videos.value.length,
    readingBooks: books.value.filter(b => b.status === '阅读中').length,
    completedBooks: books.value.filter(b => 
      b.status === '已完成' || b.status === '已解读'
    ).length,
    totalKnowledge: knowledge.value.length,
    totalSkills: skills.value.length,
    totalProblems: problems.value.length,
    totalIndex: indexNotes.value.length,
    totalQuotes: quotes.value.length,
    totalGitHub: githubRepos.value.length,
    totalAnime: anime.value.length,
    totalMovies: movies.value.length,
    totalGames: games.value.length,
    recentNotes: [],
  }))

  async function loadNotes(force = false) {
    if (loaded.value && !force) return
    
    loading.value = true
    
    try {
      const [booksData, videosData, knowledgeData, skillsData, problemsData, indexData, quotesData, githubData, animeData, moviesData, gamesData] = await Promise.all([
        booksApi.getBooks(),
        videosApi.getVideos(),
        knowledgeApi.getKnowledgeList(),
        skillsApi.getSkills(),
        problemsApi.getProblems(),
        indexApi.getIndexList(),
        quotesApi.getQuotes(),
        githubApi.getGitHubRepos(),
        animeApi.getAnime(),
        moviesApi.getMovies(),
        gamesApi.getGames(),
      ])
      
      books.value = booksData.map(booksApi.transformBookItem)
      videos.value = videosData.map(videosApi.transformVideoItem)
      knowledge.value = knowledgeData.map(knowledgeApi.transformKnowledgeItem)
      skills.value = skillsData.map(skillsApi.transformSkillItem)
      problems.value = problemsData.map(problemsApi.transformProblemItem)
      indexNotes.value = indexData.map(indexApi.transformIndexItem)
      quotes.value = quotesData.map(quotesApi.transformQuoteItem)
      githubRepos.value = githubData.map(githubApi.transformGitHubRepoItem)
      anime.value = animeData.map(animeApi.transformAnimeItem)
      movies.value = moviesData.map(moviesApi.transformMovieItem)
      games.value = gamesData.map(gamesApi.transformGameItem)
      
      loaded.value = true
    } catch (e) {
      console.error('Failed to load notes:', e)
      loaded.value = false
    } finally {
      loading.value = false
    }
  }

  async function getBookDetail(type: string, name: string): Promise<BookNote | null> {
    try {
      const data = await booksApi.getBookDetail(type, name)
      return booksApi.transformBookDetail(data)
    } catch (e) {
      console.error('Failed to load book detail:', e)
      return null
    }
  }

  async function getVideoDetail(id: string): Promise<VideoNote | null> {
    try {
      const data = await videosApi.getVideoDetail(id)
      return videosApi.transformVideoDetail(data)
    } catch (e) {
      console.error('Failed to load video detail:', e)
      return null
    }
  }

  async function getKnowledgeTree() {
    return knowledgeApi.getKnowledgeTree()
  }

  async function getKnowledgeDetail(path: string): Promise<KnowledgeNote | null> {
    try {
      const data = await knowledgeApi.getKnowledgeDetailByPath(path)
      return knowledgeApi.transformKnowledgeDetail(data)
    } catch (e) {
      console.error('Failed to load knowledge detail:', e)
      return null
    }
  }

  async function getKnowledgeDetailByQuery(category: string, topic: string, title: string): Promise<KnowledgeNote | null> {
    try {
      const data = await knowledgeApi.getKnowledgeDetailByQuery(category, topic, title)
      return knowledgeApi.transformKnowledgeDetail(data)
    } catch (e) {
      console.error('Failed to load knowledge detail:', e)
      return null
    }
  }

  async function getSkillDetail(path: string): Promise<SkillNote | null> {
    try {
      const data = await skillsApi.getSkillDetail(path)
      return skillsApi.transformSkillDetail(data)
    } catch (e) {
      console.error('Failed to load skill detail:', e)
      return null
    }
  }

  async function getProblemDetail(path: string): Promise<ProblemNote | null> {
    try {
      const data = await problemsApi.getProblemDetail(path)
      return problemsApi.transformProblemDetail(data)
    } catch (e) {
      console.error('Failed to load problem detail:', e)
      return null
    }
  }

  async function getIndexDetail(path: string): Promise<IndexNote | null> {
    try {
      const data = await indexApi.getIndexDetail(path)
      return indexApi.transformIndexDetail(data)
    } catch (e) {
      console.error('Failed to load index detail:', e)
      return null
    }
  }

  async function getQuoteDetail(id: string): Promise<QuoteNote | null> {
    try {
      const data = await quotesApi.getQuoteDetail(id)
      return quotesApi.transformQuoteDetail(data)
    } catch (e) {
      console.error('Failed to load quote detail:', e)
      return null
    }
  }

  async function getQuoteAuthors() {
    return quotesApi.getQuoteAuthors()
  }

  async function getQuoteTags() {
    return quotesApi.getQuoteTags()
  }

  async function getGitHubRepoDetail(id: string): Promise<GitHubRepoNote | null> {
    try {
      const data = await githubApi.getGitHubRepoDetail(id)
      return githubApi.transformGitHubRepoDetail(data)
    } catch (e) {
      console.error('Failed to load github repo detail:', e)
      return null
    }
  }

  async function getGitHubLanguages() {
    return githubApi.getGitHubLanguages()
  }

  async function getGitHubTags() {
    return githubApi.getGitHubTags()
  }

  // 动漫相关方法
  async function getAnimeDetail(id: string): Promise<AnimeNote | null> {
    try {
      const data = await animeApi.getAnimeDetail(id)
      return animeApi.transformAnimeDetail(data)
    } catch (e) {
      console.error('Failed to load anime detail:', e)
      return null
    }
  }

  async function getAnimeStatuses() {
    return animeApi.getAnimeStatuses()
  }

  async function getAnimeTags() {
    return animeApi.getAnimeTags()
  }

  // 电影相关方法
  async function getMovieDetail(id: string): Promise<MovieNote | null> {
    try {
      const data = await moviesApi.getMovieDetail(id)
      return moviesApi.transformMovieDetail(data)
    } catch (e) {
      console.error('Failed to load movie detail:', e)
      return null
    }
  }

  async function getMovieGenres() {
    return moviesApi.getMovieGenres()
  }

  async function getMovieTags() {
    return moviesApi.getMovieTags()
  }

  // 游戏相关方法
  async function getGameDetail(id: string): Promise<GameNote | null> {
    try {
      const data = await gamesApi.getGameDetail(id)
      return gamesApi.transformGameDetail(data)
    } catch (e) {
      console.error('Failed to load game detail:', e)
      return null
    }
  }

  async function getGamePlatforms() {
    return gamesApi.getGamePlatforms()
  }

  async function getGameStatuses() {
    return gamesApi.getGameStatuses()
  }

  async function getGameTags() {
    return gamesApi.getGameTags()
  }

  async function searchNotes(query: string) {
    return searchApi.searchNotes(query)
  }

  async function updateBookStatus(type: string, title: string, status: string): Promise<boolean> {
    return booksApi.updateBookStatus(type, title, status)
  }

  async function updateVideoStatus(id: string, status: string): Promise<boolean> {
    return videosApi.updateVideoStatus(id, status)
  }

  return {
    books,
    videos,
    knowledge,
    skills,
    problems,
    indexNotes,
    quotes,
    githubRepos,
    anime,
    movies,
    games,
    loading,
    stats,
    loadNotes,
    getBookDetail,
    getVideoDetail,
    getKnowledgeTree,
    getKnowledgeDetail,
    getKnowledgeDetailByQuery,
    getSkillDetail,
    getProblemDetail,
    getIndexDetail,
    getQuoteDetail,
    getQuoteAuthors,
    getQuoteTags,
    getGitHubRepoDetail,
    getGitHubLanguages,
    getGitHubTags,
    getAnimeDetail,
    getAnimeStatuses,
    getAnimeTags,
    getMovieDetail,
    getMovieGenres,
    getMovieTags,
    getGameDetail,
    getGamePlatforms,
    getGameStatuses,
    getGameTags,
    searchNotes,
    updateBookStatus,
    updateVideoStatus,
  }
})
