from coleta import coleta_pb2 as Coleta


def captura():
    metadado = Coleta.Metadados()
    metadado.acesso = Coleta.Metadados.FormaDeAcesso.NECESSITA_SIMULACAO_USUARIO
    metadado.extensao = Coleta.Metadados.Extensao.CSV
    metadado.estritamente_tabular = True
    metadado.tem_matricula = True
    metadado.tem_lotacao = True
    metadado.tem_cargo = True
    metadado.receita_base = Coleta.Metadados.OpcoesDetalhamento.DETALHADO
    metadado.despesas = Coleta.Metadados.OpcoesDetalhamento.DETALHADO
    metadado.outras_receitas = Coleta.Metadados.OpcoesDetalhamento.DETALHADO
    metadado.formato_consistente = True

    return metadado
