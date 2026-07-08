/**
 * API ファサード層
 *
 * アプリの全コンポーネントはこのモジュール経由で API を呼び出す。
 * VITE_API_BASE_URL が設定されていれば実 Go バックエンド (api-real)、
 * 未設定ならインメモリのモック (mock-api) を使う。モック↔実 API の
 * 切り替えはここで一元管理し、呼び出し側は切り替えを意識しない。
 */
import { isRealApiEnabled } from "@/lib/api-client";
import * as realApi from "@/lib/api-real";
import * as mockApi from "@/lib/mock-api";

const impl = isRealApiEnabled() ? realApi : mockApi;

// Users
export const getUsers = impl.getUsers;
export const getUser = impl.getUser;
export const getUsersByDojo = impl.getUsersByDojo;

// Dojos
export const getDojos = impl.getDojos;
export const getDojo = impl.getDojo;

// Practices
export const getPractices = impl.getPractices;
export const getPractice = impl.getPractice;
export const createPractice = impl.createPractice;

// Videos
export const getVideos = impl.getVideos;
export const getVideo = impl.getVideo;
export const createVideo = impl.createVideo;

// Analyses
export const getAnalyses = impl.getAnalyses;
export const getAnalysis = impl.getAnalysis;
export const getAnalysisByVideo = impl.getAnalysisByVideo;
export const analyzeVideo = impl.analyzeVideo;

// Reservations
export const getReservations = impl.getReservations;
export const getReservation = impl.getReservation;
export const createReservation = impl.createReservation;
export const deleteReservation = impl.deleteReservation;

// Exam Checklists
export const getExamChecklists = impl.getExamChecklists;
export const getExamChecklist = impl.getExamChecklist;
export const toggleChecklistItem = impl.toggleChecklistItem;

// Dashboard
export const getDashboardSummary = impl.getDashboardSummary;

// 入力・出力の型定義は実 API 実装と共有する
export type {
	AnalyzeVideoInput,
	CreatePracticeInput,
	CreateReservationInput,
	CreateVideoInput,
	DashboardSummary,
} from "@/lib/api-real";
