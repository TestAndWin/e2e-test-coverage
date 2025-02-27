export interface Area {
  id?: number;
  name?: string;
  total: number;
  passes: number;
  failures: number;
  pending: number;
  skipped: number;
  'expl-rating'?: number;
  'expl-tests'?: number;
  'first-total'?: number;
}

export interface Feature {
  id?: number;
  name?: string;
  total: number;
  passes: number;
  failures: number;
  pending: number;
  skipped: number;
  documentation: string;
  url: string;
  'business-value'?: 'low' | 'medium' | 'high';
  'first-total'?: number;
}

export interface Component {
  id?: number;
  name?: string;
  description?: string;
  'test-run'?: string;
}

export interface Test {
  id?: number;
  name?: string;
  percent?: number;
  component?: string;
  suite?: string;
  'file-name'?: string;
  'test-run'?: string;
  failures?: number;
  passes?: number;
  pending?: number;
  skipped?: number;
  total?: number;
  'failed-test-runs'?: number;
  'total-test-runs'?: number;
  url?: string;
  'area-id'?: number;
}

export interface Product {
  id?: number;
  name?: string;
  description?: string;
}

export interface User {
  id?: number;
  username?: string;
  email?: string;
  roles?: string[];
}
